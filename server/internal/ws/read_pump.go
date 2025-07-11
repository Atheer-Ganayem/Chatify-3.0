package ws

import (
	"log"
	"net"
	"time"

	"github.com/Chatify-Chat-App-in-Go-and-Next.js/server/internal/redis"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const pongWait = 60 * time.Second

func ReadPump(webSocketManager *WebSocketManager, sc *SafeConn, userID bson.ObjectID) {
	for {
		// check if rate limited
		host, _, err := net.SplitHostPort(sc.Conn.RemoteAddr().String())
		if err != nil {
			log.Println("Coudln't split host and port from client addr.")
			return
		}
		if !webSocketManager.Limiter.GetLimiter(host).Allow() {
			sc.WriteJSON(gin.H{"type": "err", "message": "Too fast."})
			continue
		}

		// Read conn & validate payload
		var payload WSPayload
		sc.Conn.SetReadDeadline(time.Now().Add(pongWait))
		err = sc.Conn.ReadJSON(&payload)
		if err != nil {
			break
		}

		conversationID, err := payload.Validate()
		if err != nil {
			sc.WriteJSON(gin.H{"type": "err", "message": err.Error()})
			continue
		}

		// check if user sent an image, if yes, validate its existience and owner in redis
		if payload.Image != "" {
			path, err := redis.GetTempImage(userID)
			if err != nil {
				sc.WriteJSON(gin.H{"type": "err", "message": "Couldn't send image."})
				continue
			} else if path != payload.Image {
				sc.WriteJSON(gin.H{"type": "err", "message": "Couldn't send image."})
				continue
			}
		}

		// saving & sending messages to other participant and ACK to client
		message, receiverID, err := payload.ProccessMessage(userID, conversationID)
		if err != nil {
			sc.WriteJSON(gin.H{"type": "err", "message": "Couldn't send message."})
			continue
		}
		go redis.DeleteKeys(userID) // cleanup redies after successfull message saving
		receiverConn := webSocketManager.GetConn(receiverID)
		if receiverConn != nil {
			if err := receiverConn.WriteJSON(gin.H{"type": "msg", "message": message}); err != nil {
				log.Printf("Failed to send to %s: %v", receiverID.Hex(), err)
			}
		}
		if err := sc.WriteJSON(gin.H{"type": "acknowledged", "message": message, "id": payload.ID}); err != nil {
			log.Printf("Failed to send to %s: %v", receiverID.Hex(), err)
		}
	}
}
