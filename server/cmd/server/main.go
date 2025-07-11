package main

import (
	"context"
	_ "crypto/tls"
	"fmt"
	"log"
	_ "net/http/httptest"
	"os"
	"time"

	routes "github.com/Chatify-Chat-App-in-Go-and-Next.js/server/internal/api"
	"github.com/Chatify-Chat-App-in-Go-and-Next.js/server/internal/db"
	"github.com/Chatify-Chat-App-in-Go-and-Next.js/server/internal/middlewares"
	"github.com/Chatify-Chat-App-in-Go-and-Next.js/server/internal/redis"
	"github.com/Chatify-Chat-App-in-Go-and-Next.js/server/internal/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}

	utils.InitAWS()
	db.Init()
	defer db.Disconnect()
	redis.Init()
	defer redis.Client.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go redis.StartSubscriber(ctx)

	server := gin.Default()
	limiter := utils.NewClientLimiter(rate.Every(750*time.Millisecond), 5)
	server.Use(middlewares.RateLimitMiddleware(limiter))

	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("FRONTEND_URL"), "http://localhost:3000"},
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
