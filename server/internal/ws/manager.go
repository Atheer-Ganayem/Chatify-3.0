package ws

import (
	"log"
	"sync"

	"github.com/Atheer-Ganayem/Chatify-3.0-backend/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type WebSocketManager struct {
	Conns     map[bson.ObjectID]*SafeConn
	ManagerMu sync.RWMutex
	Limiter   *utils.ClientLimiter
}

func NewWebSocketManager(limiter *utils.ClientLimiter) *WebSocketManager {
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
	conn, ok := wm.Conns[id]
	if ok {
		delete(wm.Conns, id)
	}
	wm.ManagerMu.Unlock()

	if ok {
		participants := conn.ParticipantsIDsCopy()
		go wm.NotifyStatus(participants, id, false)
		conn.Conn.Close()
	}
}

func (wm *WebSocketManager) GetConn(id bson.ObjectID) *SafeConn {
	wm.ManagerMu.RLock()
	defer wm.ManagerMu.RUnlock()
	return wm.Conns[id]
}

// Takes a slice of MongoDB ObjectIDs and returns only those who have an active WebSocket connection.
func (wm *WebSocketManager) FilterOnlineUsers(IDs []bson.ObjectID) []bson.ObjectID {
	wm.ManagerMu.RLock()
	defer wm.ManagerMu.RUnlock()

	var online []bson.ObjectID
	for _, id := range IDs {
		if _, ok := wm.Conns[id]; ok {
			online = append(online, id)
		}
	}
	return online
}

// Notifies all online users (from the given slice) about the online/offline status of a specific user.
func (wm *WebSocketManager) NotifyStatus(IDs []bson.ObjectID, userID bson.ObjectID, status bool) {
	onlineUsersIDs := wm.FilterOnlineUsers(IDs)
	for _, id := range onlineUsersIDs {
		go func(id bson.ObjectID) {
			conn := wm.GetConn(id)
			if conn != nil {
				if err := conn.WriteJSON(gin.H{"type": "status", "userId": userID, "online": status}); err != nil {
					log.Printf("Coudn't notify user %s of online event: %v", id, err)
				}
			}
		}(id)
	}
}
