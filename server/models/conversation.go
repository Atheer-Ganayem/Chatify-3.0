package models

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Atheer-Ganayem/Chatify-3.0-backend/db"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Conversation struct {
	ID           bson.ObjectID    `json:"_id" bson:"_id"`
	Participants [2]bson.ObjectID `json:"participants" bson:"participants"`
	LastMessage  bson.ObjectID    `json:"lastMessage,omitempty" bson:"lastMessage"`
	CreatedAt    time.Time        `json:"createdAt" bson:"createdAt"`
}

func CreateConversation(users [2]bson.ObjectID) (bson.ObjectID, int, error) {
	if len(users) != 2 {
		return bson.ObjectID{}, http.StatusBadRequest, errors.New("Users must be exactly 2.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	conversation := Conversation{ID: bson.NewObjectID(), Participants: users, CreatedAt: time.Now()}
	_, err := db.Conversations.InsertOne(ctx, conversation)

	if err != nil {
		return bson.ObjectID{}, http.StatusInternalServerError, errors.New("Something went wrong, please try again later.")
	}

	return conversation.ID, http.StatusCreated, nil
}

func FindConversation(filter bson.M) (Conversation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var conversation Conversation
	err := db.Conversations.FindOne(ctx, filter).Decode(&conversation)

	return conversation, err
}

func FindConversationByParticipants(users [2]bson.ObjectID) (Conversation, error) {
	return FindConversation(bson.M{
		"participants": bson.M{
			"$all": users,
		},
	})
}

func FindManyConversations(filter bson.M) ([]Conversation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var conversations []Conversation
	cursor, err := db.Conversations.Find(ctx, filter)
	if err != nil {
		return conversations, nil
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &conversations)

	return conversations, err
}

type PopulatedConversation struct {
	ID          bson.ObjectID `json:"_id" bson:"_id"`
	Participant UserPreview   `json:"participant"`
	LastMessage *Message      `json:"lastMessage,omitempty"`
}

type UserPreview struct {
	ID       bson.ObjectID `json:"_id" bson:"_id"`
	Name     string        `json:"name"`
	Email    string        `json:"email"`
	Avatar   string        `json:"avatar"`
	IsOnline bool          `json:"isOnline"`
}

func GetMyPopulateConversations(userID bson.ObjectID) ([]PopulatedConversation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		// Match conversations where the user is a participant
		{{Key: "$match", Value: bson.M{
			"participants": userID,
		}}},
		// Add a field: the other participant (exclude current user)
		{{Key: "$addFields", Value: bson.M{
			"otherParticipant": bson.M{
				"$filter": bson.M{
					"input": "$participants",
					"as":    "participant",
					"cond": bson.M{
						"$ne": []interface{}{"$$participant", userID},
					},
				},
			},
		}}},
		// Convert array -> object
		{{Key: "$addFields", Value: bson.M{
			"otherParticipant": bson.M{"$arrayElemAt": []interface{}{"$otherParticipant", 0}},
		}}},
		// Lookup user details
		{{Key: "$lookup", Value: bson.M{
			"from":         "users",
			"localField":   "otherParticipant",
			"foreignField": "_id",
			"as":           "participant",
		}}},
		// Convert array -> object
		{{Key: "$addFields", Value: bson.M{
			"participant": bson.M{"$arrayElemAt": []interface{}{"$participant", 0}},
		}}},
		// Lookup last message
		{{Key: "$lookup", Value: bson.M{
			"from":         "messages",
			"localField":   "lastMessage",
			"foreignField": "_id",
			"as":           "lastMessage",
		}}},
		// Convert last message array to object
		{{Key: "$addFields", Value: bson.M{
			"lastMessage": bson.M{"$arrayElemAt": []interface{}{"$lastMessage", 0}},
		}}},
		// Project only required fields
		{{Key: "$project", Value: bson.M{
			"_id":                1,
			"participant._id":    1,
			"participant.name":   1,
			"participant.email":  1,
			"participant.avatar": 1,
			"lastMessage": bson.M{
				"_id":       1,
				"sender":    1,
				"text":      1,
				"createdAt": 1,
			},
		}}},
	}

	cursor, err := db.Conversations.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var results []PopulatedConversation
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (conversation *Conversation) UpdateLastMessage(messageID bson.ObjectID) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var val interface{} = messageID

	if messageID == bson.NilObjectID {
		var lastMessage Message
		opts := options.FindOne().SetSort(bson.D{{Key: "createdAt", Value: -1}})
		err := db.Messages.FindOne(ctx, bson.M{"conversationId": conversation.ID}, opts).Decode(&lastMessage)
		if err == nil {
			val = lastMessage.ID
		} else {
			val = nil
		}
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "lastMessage", Value: val},
		}},
	}

	_, err := db.Conversations.UpdateByID(ctx, conversation.ID, update)
	if err != nil {
		log.Printf("error updating last message. conversation id: %s\n", conversation.ID)
	}
}

func GetParticipantsIDs(userID bson.ObjectID) ([]bson.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	pipeline := []bson.M{
		// Match conversations where userID is in participants array
		{
			"$match": bson.M{
				"participants": userID,
			},
		},
		// Unwind the participants array
		{
			"$unwind": "$participants",
		},
		// Filter out the current user
		{
			"$match": bson.M{
				"participants": bson.M{
					"$ne": userID,
				},
			},
		},
		// Group to get unique participants
		{
			"$group": bson.M{
				"_id": "$participants",
			},
		},
		// Project to rename _id to participant
		{
			"$project": bson.M{
				"_id":         0,
				"participant": "$_id",
			},
		},
	}

	cursor, err := db.Conversations.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []struct {
		Participant bson.ObjectID `bson:"participant"`
	}

	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	var result = make([]bson.ObjectID, len(docs))
	for i, doc := range docs {
		result[i] = doc.Participant
	}

	return result, nil
}
