package api

import (
	"log"
	"net/http"
	"time"

	"github.com/Chatify-Chat-App-in-Go-and-Next.js/server/internal/utils"
	"github.com/Chatify-Chat-App-in-Go-and-Next.js/server/internal/ws"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/time/rate"
)

var (
	webSocketManager = ws.NewWebSocketManager(utils.NewClientLimiter(rate.Every(750*time.Millisecond), 5))
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const (
	pongWait   = 60 * time.Second
	pingPeriod = 50 * time.Second
)

func connectWS(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	userID, err := bson.ObjectIDFromHex(ctx.GetString("userID"))
	if err != nil {
		log.Printf("Invalid userID: %v", err)
		return
	}

	sc := webSocketManager.ConnectUser(userID, conn)
	defer webSocketManager.DisconnectUser(userID)

	// conn.SetReadLimit(1024)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(appData string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	go ping(ticker, sc, userID)

	// notifing online and cashing user's "friends"
	err = sc.LoadParticipantsIDs(userID)
	if err != nil {
		log.Println(err.Error())
		sc.WriteJSON(gin.H{"type": "err", "message": "Couldn't load your conversations status. Try again later."})
		return
	}
	go webSocketManager.NotifyStatus(sc.ParticipantsIDsCopy(), userID, true)

	ws.ReadPump(webSocketManager, sc, userID)
}

func ping(ticker *time.Ticker, sc *ws.SafeConn, userID bson.ObjectID) {
	for range ticker.C {
		if err := sc.Ping(); err != nil {
			log.Printf("Ping failed to %s: %v", userID.Hex(), err)
			webSocketManager.DisconnectUser(userID)
			sc.Conn.Close()
			return
		}
	}
}
