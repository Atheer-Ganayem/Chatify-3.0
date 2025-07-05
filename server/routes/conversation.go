package routes

import (
	"net/http"

	"github.com/Atheer-Ganayem/Chatify-3.0-backend/models"
	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func getConversations(ctx *gin.Context) {
	userID, err := bson.ObjectIDFromHex(ctx.GetString("userID"))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Not authenticated."})
		return
	}

	conversations, err := models.GetMyPopulateConversations(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Couldn't get your conversations, please try again later."})
		return
	}

	IDs := make([]bson.ObjectID, len(conversations))
	for i, cnv := range conversations {
		IDs[i] = cnv.Participant.ID
	}
	online := webSocketManager.FilterOnlineUsers(IDs)

	ctx.JSON(http.StatusOK, gin.H{"message": "Conversation fetched successfully.",
		"conversations": conversations, "online": online})
}

func createConversation(ctx *gin.Context) {
	// Parsing & validating request body
	type CreateConversationInput struct {
		TargetUserID string `json:"targetUserID" binding:"required"`
	}

	var body CreateConversationInput
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "\"users\" field is required and must be of length 2."})
		return
	}

	// conver hex ids to object ids
	id1, err1 := bson.ObjectIDFromHex(ctx.GetString("userID"))
	id2, err2 := bson.ObjectIDFromHex(body.TargetUserID)
	if err1 != nil || err2 != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid user id."})
		return
	}

	// check if tarrget user exists
	opts := options.FindOne().SetProjection(bson.D{{Key: "_id", Value: 1}, {Key: "name", Value: 1}, {Key: "avatar", Value: 1}})
	targetUser, err := models.FindUser(bson.M{"_id": id2}, opts)
	if err == mongo.ErrNoDocuments {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Target user doesnt exist."})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong, please try again later."})
		return
	}

	// check Conversation if already exists
	userIDs := [2]bson.ObjectID{id1, id2}
	_, err = models.FindConversationByParticipants(userIDs)
	if err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Conversation already exists."})
		return
	} else if err != mongo.ErrNoDocuments {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong, please try again later."})
		return
	}

	// create Conversation
	insertedID, code, err := models.CreateConversation(userIDs)
	if err != nil {
		ctx.JSON(code, gin.H{"message": err.Error()})
		return
	}

	// Update other user of the new conversation (if connected)
	if conn := webSocketManager.GetConn(id2); conn != nil {
		conn.WriteJSON(gin.H{"type": "cnv", "user": targetUser, "cnvId": insertedID})
	}

	//response
	ctx.SecureJSON(http.StatusCreated, gin.H{"message": "Conversation has been created successfully.",
		"conversationID": insertedID})
}
