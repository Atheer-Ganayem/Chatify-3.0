package models

import (
	"context"
	"time"

	"github.com/Atheer-Ganayem/Chatify-3.0-backend/db"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const messagesLimit = 30

type Message struct {
	ID             bson.ObjectID `json:"_id" bson:"_id"`
	Sender         bson.ObjectID `json:"sender" bson:"sender"`
	ConversationID bson.ObjectID `json:"conversationId" bson:"conversationId"`
	Text           string        `json:"text" bson:"text"`
	CreatedAt      time.Time     `json:"createdAt" bson:"createdAt"`
}

func (message Message) Save() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := db.Messages.InsertOne(ctx, message)

	return err
}

func (conversation *Conversation) GetMessages(page int64) ([]Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	opts := options.Find().SetSort(bson.D{
		{Key: "createdAt", Value: -1},
	}).SetLimit(messagesLimit).SetSkip((page - 1) * messagesLimit)
	cursor, err := db.Messages.Find(ctx, bson.M{"conversationId": conversation.ID}, opts)
	if err != nil {
		return nil, err
	}

	var messages []Message
	err = cursor.All(ctx, &messages)

	return messages, err
}

func GetMessage(filter bson.M) (Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var message Message

	err := db.Messages.FindOne(ctx, filter).Decode(&message)

	return message, err
}

func (message *Message) Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.Messages.DeleteOne(ctx, bson.M{"_id": message.ID})

	return err
}
