package api

import (
	"github.com/Atheer-Ganayem/Chatify-3.0-backend/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func Register(server *gin.Engine) {
	authRoutes := server.Group("/", middlewares.IsAuth)

	{
		server.POST("/register", register)
		server.POST("/login", login)
		server.GET("/users", searchUsers)
		authRoutes.PUT("/user/name", changeUserName)
		authRoutes.PUT("/user/password", changePassword)
		authRoutes.PUT("/user/avatar", changeAvatar)
	}

	{
		authRoutes.GET("/conversations", getConversations)
		authRoutes.POST("/conversation", createConversation)
	}

	{
		authRoutes.GET("/messages/:conversationID", getMessages)
		authRoutes.DELETE("/message/:messageID", deleteMessage)
	}

	{
		authRoutes.GET("/ws", connectWS)
	}

	authRoutes.POST("/upload", uplaodHandler)

	// authRoutes.GET("/ping", func(ctx *gin.Context) {
	// 	ctx.JSON(http.StatusOK, gin.H{"message": "pong!", "Id": ctx.GetString("userID")})
	// })
}
