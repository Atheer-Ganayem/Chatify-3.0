package ws

import (
	"errors"
	"strings"
	"time"

	"github.com/Chatify-Chat-App-in-Go-and-Next.js/server-snapws/internal/models"
	"github.com/Chatify-Chat-App-in-Go-and-Next.js/server-snapws/internal/utils"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type WSPayload struct {
	ID             string `json:"id"`
	Type           string `json:"type"`
	ConversationID string `json:"conversationId"`
	Message        string `json:"message"`
	Image          string `json:"image"` // this is the value of the image path that the client received when they uploaded the image
}

func (payload *WSPayload) ProccessMessage(userID, conversationID bson.ObjectID) (models.Message, bson.ObjectID, error) {
	conversation, err := models.FindConversation(bson.M{"_id": conversationID, "participants": userID})
	if err != nil {
		return models.Message{}, bson.NewObjectID(), err
	}

	message := models.Message{ID: bson.NewObjectID(), Sender: userID, ConversationID: conversationID, Text: payload.Message, CreatedAt: time.Now()}
	if payload.Image != "" {
		message.Image = payload.Image
	}

	err = message.Save()
	if err != nil {
		return models.Message{}, bson.NewObjectID(), err
	}

	go conversation.UpdateLastMessage(message.ID)

	return message, utils.GetOtherParticipant(userID, [2]bson.ObjectID(conversation.Participants)), nil
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
