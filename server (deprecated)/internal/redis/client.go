package redis

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Chatify-Chat-App-in-Go-and-Next.js/server/internal/utils"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	DataPrefix   = "temp:image:data:"
	ExpirePrefix = "temp:image:expire:"
)

var Client *redis.Client

func Init() {
	Client = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Username: "default",
		Password: os.Getenv("REDIS_PW"),
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	result := Client.Ping(ctx)
	if result.Err() != nil {
		Client.Close()
		log.Fatalf("Redis ping failed: %v", result.Err())
	}

	fmt.Println("Redis connected successfully.")
}

// set a temp KVP in redis.
// temp:image:{userID}: {path}.
// expires in 10 minutes then cealned up.
// of temp:image:{userID} already exists, the old value will be cleaned up and replcaed by new one
func SetTempImage(path string, userID bson.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := Client.Set(ctx, ExpirePrefix+userID.Hex(), path, time.Minute*10).Err()
	if err != nil && err != redis.Nil {
		return err
	}

	prev, err := Client.SetArgs(ctx, DataPrefix+userID.Hex(), path, redis.SetArgs{
		Get: true,
		TTL: time.Minute * 12,
	}).Result()

	if err != nil {
		return err
	}
	if prev != "" {
		go utils.DeleteFile(prev)
	}

	return nil
}

func GetTempImage(userID bson.ObjectID) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	return Client.Get(ctx, DataPrefix+userID.Hex()).Result()
}

func DeleteKeys(userID bson.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	return Client.Del(ctx, DataPrefix+userID.Hex(), ExpirePrefix+userID.Hex()).Err()
}
