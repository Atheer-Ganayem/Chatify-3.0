package api

import (
	"context"
	"fmt"
	"net/http"

	snapws "github.com/Atheer-Ganayem/SnapWS"
	"github.com/Chatify-Chat-App-in-Go-and-Next.js/server-snapws/internal/models"
	"github.com/Chatify-Chat-App-in-Go-and-Next.js/server-snapws/internal/ws"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var Manager *snapws.Manager[bson.ObjectID]

func ManagerInit() {
	u := snapws.NewUpgrader(&snapws.Options{MaxMessageSize: 2048,
		ReaderMaxFragments: 5,
	})

	u.Limiter = snapws.NewRateLimiter(2, 3)
	u.Limiter.OnRateLimitHit = func(conn *snapws.Conn) error {
		conn.SendJSON(context.Background(), gin.H{"type": "err", "message": "Too fast."})
		return nil
	}

	Manager = snapws.NewManager[bson.ObjectID](u)
	Manager.OnRegister = func(conn *snapws.ManagedConn[bson.ObjectID]) {
		ids, err := models.GetParticipantsIDs(conn.Key)
		if err != nil {
			// report error to client
			conn.Close()
			return
		}
		conn.MetaData.Store("participantsIDs", ws.NewParticipantIDs(ids))

		for _, pID := range ids {
			if pConn, ok := Manager.GetConn(pID); ok {
				pConn.SendJSON(context.Background(), gin.H{"type": "status", "userId": conn.Key, "online": true})
			}
		}
	}

	Manager.OnUnregister = func(conn *snapws.ManagedConn[bson.ObjectID]) {
		val, ok := conn.MetaData.Load("participantsIDs")
		if !ok {
			return
		}

		safeIDs, ok := val.(*ws.SafeIDs)
		if !ok {
			return
		}

		safeIDs.Mu.RLock()
		defer safeIDs.Mu.RUnlock()

		for _, pID := range safeIDs.IDs {
			if pConn, ok := Manager.GetConn(pID); ok {
				pConn.SendJSON(context.Background(), gin.H{"type": "status", "userId": conn.Key, "online": false})
			}
		}
	}
}

func connectWS(ctx *gin.Context) {
	userID, err := bson.ObjectIDFromHex(ctx.GetString("userID"))
	if err != nil {
		fmt.Println("object id err")
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Not authenticated."})
		return
	}

	conn, err := Manager.Connect(userID, ctx.Writer, ctx.Request)
	if err != nil {
		fmt.Println("connect error", err)
		return
	}
	defer conn.Close()

	ws.ReadPump(Manager, conn, userID)
}

func FilterOnlineUsers(ids []bson.ObjectID) []bson.ObjectID {
	online := make([]bson.ObjectID, 0, len(ids))
	for _, id := range ids {
		if _, ok := Manager.GetConn(id); ok {
			online = append(online, id)
		}
	}

	return online
}
