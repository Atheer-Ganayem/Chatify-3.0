package utils

import (
	"sync"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type WebSocketManager struct {
	Conns map[bson.ObjectID]*websocket.Conn
	connMu          sync.RWMutex
	Limiter     *ClientLimiter
}

func NewWebSocketManager(limiter *ClientLimiter) *WebSocketManager {
	return  &WebSocketManager{
		Conns: make(map[bson.ObjectID]*websocket.Conn),
		Limiter: limiter,
	}
}

func (c *WebSocketManager) ConnectUser(userID bson.ObjectID, conn *websocket.Conn) {
	c.connMu.Lock()
	defer c.connMu.Unlock()
	c.Conns[userID] = conn
}

func (c *WebSocketManager) DisconnectUser(id bson.ObjectID) {
	c.connMu.Lock()
	defer c.connMu.Unlock()
	delete(c.Conns, id)
}

func (c *WebSocketManager) GetConn(id bson.ObjectID) *websocket.Conn {
	c.connMu.RLock()
	defer c.connMu.RUnlock()
	return c.Conns[id]
}
