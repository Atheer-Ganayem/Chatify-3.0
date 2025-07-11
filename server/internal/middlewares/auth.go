package middlewares

import (
	"net/http"
	"strings"

	"github.com/Chatify-Chat-App-in-Go-and-Next.js/server/internal/models"
	"github.com/Chatify-Chat-App-in-Go-and-Next.js/server/internal/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func IsAuth(ctx *gin.Context) {
	var token string

	if ctx.Request.Header.Get("Upgrade") == "websocket" {
		token = ctx.Query("token")
		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token query."})
			return
		}
	} else {
		authHeader := strings.Split(ctx.GetHeader("Authorization"), " ")

		if len(authHeader) != 2 || authHeader[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid auth header."})
			return
		}
		token = authHeader[1]
	}

	userHexID, err := utils.VerifyToken(token)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Authentication is required."})
		return
	}

	userObjectID, err := bson.ObjectIDFromHex(userHexID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Authentication is required."})
		return
	}
	exists, err := models.UserExists(bson.M{"_id": userObjectID})
	if !exists || err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Authentication is required."})
		return
	}

	ctx.Set("userID", userHexID)
	ctx.Next()
}
