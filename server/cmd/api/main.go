package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/victoryus84/gorders/internal/config"
	"github.com/victoryus84/gorders/internal/database"
	"github.com/victoryus84/gorders/internal/handler"
	"github.com/victoryus84/gorders/internal/kafka"
	"github.com/victoryus84/gorders/internal/logger"
	"github.com/victoryus84/gorders/internal/middleware"
	"github.com/victoryus84/gorders/internal/repository"
	"github.com/victoryus84/gorders/internal/router"
	"github.com/victoryus84/gorders/internal/service"
)

// Build flags set during compilation
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildTime = "unknown"
)

func main() {
	_ = godotenv.Load()
	// Load configuration (singleton)
	cfg := config.Load()

	Version = cfg.Version
	Commit = cfg.Commit

	if cfg.AppEnv == "production" {
		logger.Init("info")
		logger.LogInfo("⚠️ RUNNING IN PRODUCTION MODE",
			logger.String("version", Version), // Uite ce curat arată!
			logger.String("commit", Commit),
		)
	} else {
		logger.Init("debug")
		logger.LogInfo("🛠️ Running in Development mode",
			logger.String("version", Version),
		)
	}

	// Pornim Producer-ul cu adresa din .env
	kp := kafka.NewProducer(cfg.KafkaAddr)
	defer kp.Close()

	// Initialize structured logging
	logger.Init(cfg.LogLevel)
	defer logger.Logger.Sync()

	// Print startup banner
	printBanner()

	// Connect to database
	db := database.Connect(cfg)

	// Create repository
	rep := repository.NewRepository(db)

	// Create services
	svc_user := service.NewUserService(rep, cfg)
	svc_client := service.NewClientService(rep, cfg, kp)
	svc_contract := service.NewContractService(rep, rep, cfg, kp)

	logger.LogInfo("✅ All services initialized")

	// Setup Gin router
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Apply global middleware in order
	r.Use(middleware.RequestLogging()) // Logging with trace ID
	r.Use(middleware.PanicRecovery())  // Panic recovery
	r.Use(middleware.CORS())           // CORS headers
	r.Use(middleware.RateLimit())      // Rate limiting

	// Core handlers
	hdl_core := handler.NewCoreHandler(db, Version, Commit)
	hdl_user := handler.NewUserHandler(svc_user)
	hdl_client := handler.NewClientHandler(svc_client)
	hdl_contract := handler.NewContractHandler(svc_contract)

	allHandlers := &handler.Handlers{
		Core:     hdl_core,
		User:     hdl_user,
		Client:   hdl_client,
		Contract: hdl_contract,
	}
	// Setup API routes
	router.SetupRoutes(r, allHandlers)

	logger.LogInfo("🎯 Server starting",
		logger.String("port", "8080"),
		logger.String("env", cfg.AppEnv),
		logger.String("version", Version),
		logger.String("commit", Commit),
		logger.String("buildTime", BuildTime),
	)

	// Start server
	if err := r.Run(":8080"); err != nil {
		// Asta va scrie log-ul frumos în JSON (cu Zap) ȘI va opri serverul!
		logger.LogFatal("Server failed to start", err)
	}
}

// printBanner prints startup banner
func printBanner() {
	// 1. Folosim backticks (`) pentru a scrie pe mai multe rânduri
	// FĂRĂ să facem concatenări urâte. Sprintf rezolvă formatarea cu variabilele tale.
	banner := fmt.Sprintf(`
╔════════════════════════════════════════╗
║     🚀 GOrders Backend Server 🚀       ║
╠════════════════════════════════════════╣
║  Version:  %-26s  ║
║  Commit:   %-26s  ║
║  Built:    %-26s  ║
╚════════════════════════════════════════╝
`, Version, Commit, BuildTime)

	// 2. Acum avem un singur log. În development, va apărea cu verde "INFO" deasupra.
	// Pe server, va fi un singur JSON curat.
	logger.LogInfo(banner)
}
