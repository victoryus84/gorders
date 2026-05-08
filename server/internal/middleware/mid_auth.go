package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/victoryus84/gorders/internal/config"
	"github.com/victoryus84/gorders/internal/logger"
	"go.uber.org/zap"
)

// AuthMiddleware validates JWT tokens
func AuthJWT() gin.HandlerFunc {
	cfg := config.Load() // Singleton - called once, reused

	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			logger.LogWarn("Invalid token attempt",
				zap.String("trace_id", c.GetString("trace_id")),
				zap.Error(err),
			)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Extract user info from claims
		userID, _ := claims["user_id"].(float64)
		role, _ := claims["role"].(string)

		c.Set("user_id", uint(userID))
		c.Set("role", role)
		c.Next()
	}
}
