package ws

import (
	"context"
	"fmt"
	"log"

	snapws "github.com/Atheer-Ganayem/SnapWS"
	"github.com/Chatify-Chat-App-in-Go-and-Next.js/server-snapws/internal/redis"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func ReadPump(manager *snapws.Manager[bson.ObjectID], conn *snapws.ManagedConn[bson.ObjectID], userID bson.ObjectID) {
	for {
		// Read conn & validate payload
		var payload WSPayload
		err := conn.ReadJSON(&payload)
		if snapws.IsFatalErr(err) {
			/// here it error
			fmt.Println("reader fatal", err)
			return
		} else if err != nil {
			fmt.Println(err)
			// must report to client
			continue
		}

		conversationID, err := payload.Validate()
		if err != nil {
			fmt.Println(payload)
			conn.SendJSON(context.Background(), gin.H{"type": "err", "message": err.Error()})
			continue
		}

		// check if user sent an image, if yes, validate its existience and owner in redis
		if payload.Image != "" {
			path, err := redis.GetTempImage(userID)
			if err != nil {
				conn.SendJSON(context.Background(), gin.H{"type": "err", "message": "Couldn't send image."})
				continue
			} else if path != payload.Image {
				conn.SendJSON(context.Background(), gin.H{"type": "err", "message": "Couldn't send image."})
				continue
			}
		}

		// saving & sending messages to other participant and ACK to client
		message, receiverID, err := payload.ProccessMessage(userID, conversationID)
		if err != nil {
			conn.SendJSON(context.Background(), gin.H{"type": "err", "message": "Couldn't send message."})
			continue
		}
		go redis.DeleteKeys(userID) // cleanup redies after successfull message saving

		receiverConn, ok := manager.GetConn(receiverID)
		if receiverConn != nil && ok {
			if err := receiverConn.SendJSON(context.Background(), gin.H{"type": "msg", "message": message}); err != nil {
				log.Printf("Failed to send to %s: %v", receiverID.Hex(), err)
			}
		}
		if err := conn.SendJSON(context.Background(), gin.H{"type": "acknowledged", "message": message, "id": payload.ID}); err != nil {
			log.Printf("Failed to send to %s: %v", receiverID.Hex(), err)
		}
	}
}
