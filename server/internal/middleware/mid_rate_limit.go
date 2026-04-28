package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/victoryus84/gorders/internal/config"
)

// RateLimiter implements token bucket rate limiting
type RateLimiter struct {
	clients map[string]*clientLimiter
	mu      sync.Mutex
	config  *config.Config
}

type clientLimiter struct {
	tokens    int
	lastReset time.Time
}

var rateLimiter *RateLimiter

func init() {
	cfg := config.Load()
	rateLimiter = &RateLimiter{
		clients: make(map[string]*clientLimiter),
		config:  cfg,
	}

	// Cleanup goroutine - remove old entries every minute
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for range ticker.C {
			rateLimiter.cleanup()
		}
	}()
}

// RateLimit middleware enforces rate limiting per IP
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !rateLimiter.Allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Too many requests",
				"message": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Allow checks if request from IP is allowed
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	limiter, exists := rl.clients[ip]

	if !exists {
		rl.clients[ip] = &clientLimiter{
			tokens:    rl.config.RateLimitRequests - 1,
			lastReset: now,
		}
		return true
	}

	// Check if window has passed
	if now.Sub(limiter.lastReset) > time.Duration(rl.config.RateLimitWindow)*time.Second {
		limiter.tokens = rl.config.RateLimitRequests - 1
		limiter.lastReset = now
		return true
	}

	// Check if tokens available
	if limiter.tokens > 0 {
		limiter.tokens--
		return true
	}

	return false
}

// cleanup removes old client entries
func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for ip, limiter := range rl.clients {
		if now.Sub(limiter.lastReset) > 5*time.Minute {
			delete(rl.clients, ip)
		}
	}
}
