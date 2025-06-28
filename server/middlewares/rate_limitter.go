package middlewares

import (
	"net"
	"net/http"

	"github.com/Atheer-Ganayem/Chatify-3.0-backend/utils"
	"github.com/gin-gonic/gin"
)

func RateLimitMiddleware(cl *utils.ClientLimiter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := getClientIP(ctx)
		limiter := cl.GetLimiter(ip)

		if !limiter.Allow() {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"message": "Too many requests"})
			return
		}

		ctx.Next()
	}
}

func getClientIP(ctx *gin.Context) string {
	ip := ctx.ClientIP()

	host, _, err := net.SplitHostPort(ip)
	if err == nil {
		return host
	}

	return ip
}
