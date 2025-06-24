package middlewares

import (
	"net/http"
	"strings"

	"github.com/Atheer-Ganayem/Chatify-3.0-backend/models"
	"github.com/Atheer-Ganayem/Chatify-3.0-backend/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func IsAuth(ctx *gin.Context) {
	found := true
	var token string
	if ctx.Request.Header.Get("Upgrade") == "websocket" {
		token = ctx.Query("token")
		if token == "" {
			found = false
		}
	} else {
		authHeader := strings.Split(ctx.GetHeader("Authorization"), " ")
		if len(authHeader) != 2 || authHeader[0] != "Bearer" {
			found = false
		}
		token = authHeader[1]
	}

	if !found {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid auth header."})
		return
	}

	userHexID, err := utils.VerifyToken(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Authentication is required."})
		return
	}

	userObjectID, err := bson.ObjectIDFromHex(userHexID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Authentication is required."})
		return
	}
	exists, err := models.UserExists(bson.M{"_id": userObjectID})
	if !exists || err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Authentication is required."})
		return
	}

	ctx.Set("userID", userHexID)
	ctx.Next()
}
