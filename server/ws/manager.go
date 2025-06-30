package ws

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Atheer-Ganayem/Chatify-3.0-backend/models"
	"github.com/Atheer-Ganayem/Chatify-3.0-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const writeWait = 8 * time.Second

type SafeConn struct {
	Conn            *websocket.Conn
	ParticipantsIDs []bson.ObjectID
	mu              sync.RWMutex
}

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

func (sc *SafeConn) WriteJSON(payload gin.H) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.Conn.SetWriteDeadline(time.Now().Add(writeWait))
	return sc.Conn.WriteJSON(payload)
}

func (sc *SafeConn) Ping() error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.Conn.SetWriteDeadline(time.Now().Add(writeWait))
	return sc.Conn.WriteMessage(websocket.PingMessage, nil)
}

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

func (sc *SafeConn) LoadParticipantsIDs(userID bson.ObjectID) error {
	IDs, err := models.GetParticipantsIDs(userID)
	if err != nil {
		return fmt.Errorf("failed to load participants for user %s: %w", userID.Hex(), err)
	}

	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.ParticipantsIDs = IDs

	return nil
}

func (sc *SafeConn) ParticipantsIDsCopy() []bson.ObjectID {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	idsCopy := make([]bson.ObjectID, len(sc.ParticipantsIDs))
	copy(idsCopy, sc.ParticipantsIDs)

	return idsCopy
}
