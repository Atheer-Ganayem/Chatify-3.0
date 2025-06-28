package routes

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/Atheer-Ganayem/Chatify-3.0-backend/models"
	"github.com/Atheer-Ganayem/Chatify-3.0-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/time/rate"
)

var (
	webSocketManager = utils.NewWebSocketManager(utils.NewClientLimiter(rate.Every(750*time.Millisecond), 5))
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const (
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10 // 54 seconds
	writeWait  = 10 * time.Second
)

func connectWS(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	userObjectID, err := bson.ObjectIDFromHex(ctx.GetString("userID"))
	if err != nil {
		log.Printf("Invalid userID: %v", err)
		return
	}

	webSocketManager.ConnectUser(userObjectID, conn)
	defer webSocketManager.DisconnectUser(userObjectID)

	// conn.SetReadLimit(1024)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(appData string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	go ping(ticker, conn, userObjectID)

	for {
		// check if rate limited
		host, _, err := net.SplitHostPort(conn.RemoteAddr().String())
		if err != nil {
			log.Println("Coudln't split host and port from client addr.")
			return
		}
		if !webSocketManager.Limiter.GetLimiter(host).Allow() {
			conn.WriteJSON(gin.H{"type": "err", "message": "Too fast."})
			continue
		}

		// validate payload
		var payload models.WSPayload
		err = conn.ReadJSON(&payload)
		if err != nil {
			log.Printf("Read error from %s: %v", userObjectID.Hex(), err)
			break
		}

		conversationID, err := payload.Validate()
		if err != nil {
			conn.WriteJSON(gin.H{"type": "err", "message": err.Error()})
			continue
		}

		// saving & sending messages to other participant and ACK to client
		message, receiverID, err := payload.SaveMessage(userObjectID, conversationID)
		if err != nil {
			conn.WriteJSON(gin.H{"type": "err", "message": "Couldn't send message."})
			continue
		}
		receiverConn := webSocketManager.GetConn(receiverID)
		if receiverConn != nil {
			if err := receiverConn.WriteJSON(gin.H{"type": "msg", "message": message}); err != nil {
				log.Printf("Failed to send to %s: %v", receiverID.Hex(), err)
			}
		}
		if err := conn.WriteJSON(gin.H{"type": "acknowledged", "message": message, "id": payload.ID}); err != nil {
			log.Printf("Failed to send to %s: %v", receiverID.Hex(), err)
		}
	}
}

func ping(ticker *time.Ticker, conn *websocket.Conn, userObjectID bson.ObjectID) {
	for range ticker.C {
		conn.SetWriteDeadline(time.Now().Add(writeWait))
		if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			log.Printf("Ping failed to %s: %v", userObjectID.Hex(), err)
			conn.Close()
			return
		}
	}
}
