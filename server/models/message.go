package models

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/Atheer-Ganayem/Chatify-3.0-backend/db"
	"github.com/Atheer-Ganayem/Chatify-3.0-backend/utils"

	"github.com/google/uuid"
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

type WSPayload struct {
	ID             string `json:"id"`
	Type           string `json:"type"`
	ConversationID string `json:"conversationId"`
	Message        string `json:"message"`
}

func (payload *WSPayload) SaveMessage(userID, conversationID bson.ObjectID) (Message, bson.ObjectID, error) {
	conversation, err := FindConversation(bson.M{"_id": conversationID, "participants": userID})
	if err != nil {
		return Message{}, bson.NewObjectID(), err
	}

	message := Message{ID: bson.NewObjectID(), Sender: userID, ConversationID: conversationID, Text: payload.Message, CreatedAt: time.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err = db.Messages.InsertOne(ctx, message)
	if err != nil {
		return Message{}, bson.NewObjectID(), err
	}

	go conversation.UpdateLastMessage(message.ID)

	return message, utils.GetOtherParticipant(userID, [2]bson.ObjectID(conversation.Participants)), nil
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

func (payload *WSPayload) Validate() (bson.ObjectID, error) {
	payload.Message = strings.TrimSpace(payload.Message)
	if _, err := uuid.Parse(payload.ID); err != nil {
		return bson.NilObjectID, errors.New("Invalid request ID. Must be a valid UUID.")
	} else if payload.Type != "msg" {
		return bson.NilObjectID, errors.New("Invalid type.")
	} else if payload.Message == "" {
		return bson.NilObjectID, errors.New("Message cannot be empty.")
	}

	conversationID, err := bson.ObjectIDFromHex(payload.ConversationID)
	if err != nil {
		return bson.NilObjectID, errors.New("Invalid conversation ID.")
	}

	return conversationID, nil
}
