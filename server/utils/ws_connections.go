package utils

import (
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const writeWait = 8 * time.Second

type SafeConn struct {
	Conn *websocket.Conn
	mu   sync.Mutex
}

type WebSocketManager struct {
	Conns     map[bson.ObjectID]*SafeConn
	ManagerMu sync.RWMutex
	Limiter   *ClientLimiter
}

func NewWebSocketManager(limiter *ClientLimiter) *WebSocketManager {
	return &WebSocketManager{
		Conns:   make(map[bson.ObjectID]*SafeConn),
		Limiter: limiter,
	}
}

func (wm *WebSocketManager) ConnectUser(userID bson.ObjectID, conn *websocket.Conn) *SafeConn {
	wm.ManagerMu.Lock()
	defer wm.ManagerMu.Unlock()

	if old, ok := wm.Conns[userID]; ok {
		old.Conn.Close()
	}

	sc := &SafeConn{Conn: conn}
	wm.Conns[userID] = sc
	return sc
}

func (wm *WebSocketManager) DisconnectUser(id bson.ObjectID) {
	wm.ManagerMu.Lock()
	defer wm.ManagerMu.Unlock()

	if conn, ok := wm.Conns[id]; ok {
		conn.Conn.Close()
	}
	delete(wm.Conns, id)
}

func (wm *WebSocketManager) GetConn(id bson.ObjectID) *SafeConn {
	wm.ManagerMu.RLock()
	defer wm.ManagerMu.RUnlock()
	return wm.Conns[id]
}

func (sc *SafeConn) WriteJSON(payload gin.H) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.Conn.SetWriteDeadline(time.Now().Add(writeWait))
	return sc.Conn.WriteJSON(payload)
}

func (sc *SafeConn) Ping() error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	log.Println("Pinging...")

	sc.Conn.SetWriteDeadline(time.Now().Add(writeWait))
	return sc.Conn.WriteMessage(websocket.PingMessage, nil)
}
