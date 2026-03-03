package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
}

type visitor struct {
	count    int
	resetAt  time.Time
}

var limiter = &rateLimiter{
	visitors: make(map[string]*visitor),
}

func init() {
	// Clean up stale entries every 5 minutes.
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			limiter.mu.Lock()
			now := time.Now()
			for ip, v := range limiter.visitors {
				if now.After(v.resetAt) {
					delete(limiter.visitors, ip)
				}
			}
			limiter.mu.Unlock()
		}
	}()
}

// RateLimit allows max requests per minute per IP.
func RateLimit(maxPerMinute int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		limiter.mu.Lock()
		v, exists := limiter.visitors[ip]
		if !exists || now.After(v.resetAt) {
			limiter.visitors[ip] = &visitor{count: 1, resetAt: now.Add(time.Minute)}
			limiter.mu.Unlock()
			c.Next()
			return
		}

		v.count++
		if v.count > maxPerMinute {
			limiter.mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded, try again later"})
			return
		}
		limiter.mu.Unlock()
		c.Next()
	}
}
