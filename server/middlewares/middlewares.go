package middlewares

import (
	"net/http"

	"github.com/Atheer-Ganayem/Chatify-3.0-backend/models"
	"github.com/Atheer-Ganayem/Chatify-3.0-backend/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func IsAuth(ctx *gin.Context) {
	authCookie, err := ctx.Cookie("next-auth.session-token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "No auth cookie."})
		return
	}

	userHexID, err := utils.VerifyToken(authCookie)
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
