package redis

import (
	"context"
	"log"
	"strings"

	"github.com/Chatify-Chat-App-in-Go-and-Next.js/server-snapws/internal/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func StartSubscriber(ctx context.Context) {
	pubsub := Client.PSubscribe(ctx, "__keyevent@0__:expired")
	ch := pubsub.Channel()

	go func() {
		defer pubsub.Close()
		for {
			select {
			case <-ctx.Done():
				log.Println("Redis subscriber shutting down")
				return
			case msg, ok := <-ch:
				if !ok {
					log.Println("PubSub channel closed")
					return
				}

				key := msg.Payload
				if strings.HasPrefix(key, ExpirePrefix) {
					go handleExpiry(key)
				}
			}
		}
	}()
}

func handleExpiry(key string) {
	hexID := strings.TrimPrefix(key, ExpirePrefix)
	userID, err := bson.ObjectIDFromHex(hexID)
	if err != nil {
		log.Printf("Coudn't convert redis key to object id: %s\n", err.Error())
		return
	}

	path, err := GetTempImage(userID)
	if err != nil {
		log.Printf("Coudn't get image path from redis: %s\n", err.Error())
		return
	}

	go utils.DeleteFile(path)
}
