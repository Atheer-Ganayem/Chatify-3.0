package ws

import (
	"fmt"
	"sync"
	"time"

	"github.com/Atheer-Ganayem/Chatify-3.0-backend/models"
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

func (sc *SafeConn) ParticipantsIDsCopy() []bson.ObjectID {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	idsCopy := make([]bson.ObjectID, len(sc.ParticipantsIDs))
	copy(idsCopy, sc.ParticipantsIDs)

	return idsCopy
}

// Loads participant IDs for the user from MongoDB and stores them in the SafeConn.
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
