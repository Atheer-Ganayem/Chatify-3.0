package api

import (
	"net/http"

	"github.com/Chatify-Chat-App-in-Go-and-Next.js/server/internal/redis"
	"github.com/Chatify-Chat-App-in-Go-and-Next.js/server/internal/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func uplaodHandler(ctx *gin.Context) {
	userID, err := bson.ObjectIDFromHex(ctx.GetString("userID"))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Not authenticated."})
		return
	}

	path, code, err := utils.ExtractFileAndUpload(ctx.Request, "image")
	if err != nil {
		ctx.JSON(code, gin.H{"message": err.Error()})
		return
	}

	err = redis.SetTempImage(path, userID)
	if err != nil {
		go utils.DeleteFile(path)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Couldn't upload image, please try again later."})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Image uploaded successfully.", "path": path})
}

func deleteHandler(ctx *gin.Context) {
	userID, err := bson.ObjectIDFromHex(ctx.GetString("userID"))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Not authenticated."})
		return
	}

	path, err := redis.GetTempImage(userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Image not found."})
		return
	}

	go utils.DeleteFile(path)

	ctx.SecureJSON(http.StatusOK, gin.H{"message": "Deleting image."})
}
