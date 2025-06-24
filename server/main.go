package main

import (
	_ "crypto/tls"
	"fmt"
	"log"
	_ "net/http/httptest"
	"os"

	"github.com/Atheer-Ganayem/Chatify-3.0-backend/db"
	"github.com/Atheer-Ganayem/Chatify-3.0-backend/routes"
	"github.com/Atheer-Ganayem/Chatify-3.0-backend/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}

	utils.InitAWS()
	db.Init()
	defer db.Disconnect()

	server := gin.Default()

	server.Use(cors.New(cors.Config{
		// AllowAllOrigins: true,
		AllowOrigins:     []string{os.Getenv("FRONTEND_URL")},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "Accept", "Origin"},
		AllowCredentials: true,
	}))

	routes.Register(server)

	port := os.Getenv("PORT")
	if port == "" {
		port = "localhost:8080"
	} else {
		port = fmt.Sprintf(":%s", port)
	}
	server.Run(port)
}
