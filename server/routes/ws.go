package routes

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Atheer-Ganayem/Chatify-3.0-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type connections map[bson.ObjectID]*websocket.Conn

var (
	ConnectedUsers = make(connections)
	connMu         sync.RWMutex
)

func (c *connections) connectUser(userID bson.ObjectID, conn *websocket.Conn) {
	connMu.Lock()
	defer connMu.Unlock()
	(*c)[userID] = conn
}

func (c *connections) disconnectUser(id bson.ObjectID) {
	connMu.Lock()
	defer connMu.Unlock()
	delete((*c), id)
}

func (c *connections) getConn(id bson.ObjectID) *websocket.Conn {
	connMu.RLock()
	defer connMu.RUnlock()
	return (*c)[id]
}

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

	ConnectedUsers.connectUser(userObjectID, conn)
	defer ConnectedUsers.disconnectUser(userObjectID)

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
		var payload models.WSPayload
		err := conn.ReadJSON(&payload)
		if err != nil {
			log.Printf("Read error from %s: %v", userObjectID.Hex(), err)
			break
		}

		if payload.Type == "ping" {
			if err := conn.WriteJSON(gin.H{"type": "pong"}); err != nil {
				log.Printf("Failed to send pong to %s: %v", userObjectID.Hex(), err)
				break
			}
			continue
		}

		if payload.Message == "" || payload.ConversationID == "" || payload.ID == "" {
			conn.WriteJSON(gin.H{"type": "err", "message": "Message, conversation ID or request ID is missing."})
			continue
		}

		conversationID, err := bson.ObjectIDFromHex(payload.ConversationID)
		if err != nil {
			conn.WriteJSON(gin.H{"type": "err", "message": "Invalid conversation ID."})
			log.Printf("Invalid conversation ID: %v", err)
			continue
		}
		message, receiverID, err := payload.SaveMessage(userObjectID, conversationID)
		if err != nil {
			conn.WriteJSON(gin.H{"type": "err", "message": "Couldn't send message."})
			log.Printf("SaveMessage failed: %v", err)
			continue
		}
		receiverConn := ConnectedUsers.getConn(receiverID)
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
