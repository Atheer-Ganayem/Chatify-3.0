package middlewares

import (
	"net"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type clientLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.Mutex
	rate     rate.Limit
	burst    int
}

func NewClient(r rate.Limit, b int) *clientLimiter {
	return &clientLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     r,
		burst:    b,
	}
}

func (cl *clientLimiter) getLimiter(ip string) *rate.Limiter {
	cl.mu.Lock()
	defer cl.mu.Unlock()

	limiter, exists := cl.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(cl.rate, cl.burst)
		cl.limiters[ip] = limiter
	}

	return limiter
}

func getClientIP(ctx *gin.Context) string {
	ip := ctx.ClientIP()

	host, _, err := net.SplitHostPort(ip)
	if err == nil {
		return host
	}

	return ip
}

func RateLimitMiddleware(cl *clientLimiter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := getClientIP(ctx)
		limiter := cl.getLimiter(ip)

		if !limiter.Allow() {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"message": "Too many requests"})
			return 
		}

		ctx.Next()
	}
}
