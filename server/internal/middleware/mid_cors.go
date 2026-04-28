package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/victoryus84/gorders/internal/config"
)

// CORS configures Cross-Origin Resource Sharing
func CORS() gin.HandlerFunc {
	cfg := config.Load()

	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", cfg.CORSAllowedOrigins)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", cfg.CORSAllowedMethods)
		c.Writer.Header().Set("Access-Control-Allow-Headers", cfg.CORSAllowedHeaders)

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
