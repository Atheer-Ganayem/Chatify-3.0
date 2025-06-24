package routes

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Atheer-Ganayem/Chatify-3.0-backend/models"
	"github.com/Atheer-Ganayem/Chatify-3.0-backend/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func getMessages(ctx *gin.Context) {
	conversationHexID, _ := ctx.Params.Get("conversationID")
	conversationObjectID, err := bson.ObjectIDFromHex(conversationHexID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "conversationID param is required."})
		return
	}

	userObjectID, err := bson.ObjectIDFromHex(ctx.GetString("userID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid or missing user ID."})
		return
	}

	conversation, err := models.FindConversation(bson.M{"_id": conversationObjectID, "participants": userObjectID})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Conversation not found."})
		return
	}

	page, err := strconv.ParseInt(ctx.Query("page"), 10, 64)
	if err != nil || page <= 0 {
		page = 1
	}

	messages, err := conversation.GetMessages(page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Couldn't fetch messages."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Messages fetched successfully", "messages": messages})
}

func deleteMessage(ctx *gin.Context) {
	messageHexID, _ := ctx.Params.Get("messageID")
	messageID, err := bson.ObjectIDFromHex(messageHexID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid message id."})
		return
	}

	userID, err := bson.ObjectIDFromHex(ctx.GetString("userID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid or missing user ID."})
		return
	}

	message, err := models.GetMessage(bson.M{"_id": messageID, "sender": userID})
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Message not found (1)."})
		return
	}

	err = message.Delete()

	if err == mongo.ErrNoDocuments {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Message not found (2)."})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Coudln't delete message, please try again later."})
		return
	}

	conversation, err := models.FindConversation(bson.M{"participants": userID})
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "Message has been deleted successuflly. But conversation not found."})
		return
	}

	go conversation.UpdateLastMessage(bson.NilObjectID)

	otherUserID := utils.GetOtherParticipant(userID, conversation.Participants)
	if conn := ConnectedUsers.getConn(otherUserID); conn != nil {
		if err = conn.WriteJSON(gin.H{"type": "delete", "messageId": messageID}); err != nil {
			log.Println("Coudn't send delete update via ws to other user.")
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Message has been deleted successuflly."})
}
