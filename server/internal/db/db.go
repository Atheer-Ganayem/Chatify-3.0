package db

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	client        *mongo.Client
	DB            *mongo.Database
	Users         *mongo.Collection
	Conversations *mongo.Collection
	Messages      *mongo.Collection
)

func Init() {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		panic("You must set your 'MONGODB_URI' environment variable.")
	}

	opts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(opts)
	if err != nil {
		panic(err)
	}

	DB = client.Database("Chatify-3")
	Users = DB.Collection("users")
	Conversations = DB.Collection("conversations")
	Messages = DB.Collection("messages")

	log.Println("DB connected!")
}

func Disconnect() {
	if err := client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}

}
