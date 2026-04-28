package logger

import (
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

// Init initializes the global logger
func Init(logLevel string) {
	var config zap.Config

	if os.Getenv("APP_ENV") == "production" {
		config = zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(parseLogLevel(logLevel))
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	var err error
	Logger, err = config.Build()
	if err != nil {
		panic(err)
	}
	defer Logger.Sync()
}

// parseLogLevel converts string to zapcore level
func parseLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// LogRequest logs HTTP request details
func LogRequest(c *gin.Context) {
	Logger.Info("HTTP Request",
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.String("query", c.Request.URL.RawQuery),
		zap.String("ip", c.ClientIP()),
		zap.String("trace_id", c.GetString("trace_id")),
	)
}

// LogResponse logs HTTP response details
func LogResponse(c *gin.Context, latency float64) {
	Logger.Info("HTTP Response",
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.Int("status", c.Writer.Status()),
		zap.Float64("latency_ms", latency),
		zap.String("trace_id", c.GetString("trace_id")),
	)
}

// LogError logs error with context
func LogError(message string, err error, fields ...zap.Field) {
	Logger.Error(message,
		append([]zap.Field{zap.Error(err)}, fields...)...,
	)
}

// LogInfo logs informational message
func LogInfo(message string, fields ...zap.Field) {
	Logger.Info(message, fields...)
}

// LogWarn logs warning message
func LogWarn(message string, fields ...zap.Field) {
	Logger.Warn(message, fields...)
}

// LogDebug logs debug message
func LogDebug(message string, fields ...zap.Field) {
	Logger.Debug(message, fields...)
}
