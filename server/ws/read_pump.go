package ws

import (
	"log"
	"net"
	"time"

	"github.com/Atheer-Ganayem/Chatify-3.0-backend/models"
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
		var payload models.WSPayload
		sc.Conn.SetReadDeadline(time.Now().Add(pongWait))
		err = sc.Conn.ReadJSON(&payload)
		if err != nil {
			log.Printf("Read error from %s: %v", userID.Hex(), err)
			break
		}

		conversationID, err := payload.Validate()
		if err != nil {
			sc.WriteJSON(gin.H{"type": "err", "message": err.Error()})
			continue
		}

		// saving & sending messages to other participant and ACK to client
		message, receiverID, err := payload.SaveMessage(userID, conversationID)
		if err != nil {
			sc.WriteJSON(gin.H{"type": "err", "message": "Couldn't send message."})
			continue
		}
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
