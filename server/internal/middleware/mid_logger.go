package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/victoryus84/gorders/internal/logger"
	"go.uber.org/zap"
)

// RequestLogging logs all HTTP requests and responses with trace ID
func RequestLogging() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate trace ID for request tracking
		traceID := uuid.New().String()
		c.Set("trace_id", traceID)

		startTime := time.Now()

		// Log incoming request
		logger.LogInfo("Incoming request",
			zap.String("trace_id", traceID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("ip", c.ClientIP()),
		)

		c.Next()

		// Log response
		latency := time.Since(startTime).Milliseconds()
		logger.LogInfo("Request completed",
			zap.String("trace_id", traceID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Int64("latency_ms", latency),
		)
	}
}

// ErrorResponse is used for error logging
func ErrorResponse(c *gin.Context, statusCode int, message string, traceID string) {
	c.JSON(statusCode, gin.H{
		"error":    http.StatusText(statusCode),
		"message":  message,
		"trace_id": traceID,
	})
}
