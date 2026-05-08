package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/victoryus84/gorders/internal/logger"
	"go.uber.org/zap"
)

// PanicRecovery recovers from panics and logs them
func PanicRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				traceID := c.GetString("trace_id")
				logger.LogError("Panic recovered",
					nil,
					zap.String("trace_id", traceID),
					zap.Any("panic", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
				)

				c.JSON(http.StatusInternalServerError, gin.H{
					"error":    "Internal Server Error",
					"message":  "An unexpected error occurred",
					"trace_id": traceID,
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
