package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/victoryus84/gorders/server/internal/config"
	"github.com/victoryus84/gorders/server/internal/database"
	"github.com/victoryus84/gorders/server/internal/handler"
	"github.com/victoryus84/gorders/server/internal/kafka"
	"github.com/victoryus84/gorders/server/internal/logger"
	"github.com/victoryus84/gorders/server/internal/middleware"
	"github.com/victoryus84/gorders/server/internal/repository"
	"github.com/victoryus84/gorders/server/internal/router"
	"github.com/victoryus84/gorders/server/internal/service"
	"go.uber.org/zap"
)

// Build flags set during compilation
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildTime = "unknown"
)

func main() {
	// Load configuration (singleton)
	cfg := config.Load()

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
	repo := repository.NewRepository(db)

	// Create services
	svc_user := service.NewUserService(repo, cfg) 
	svc_client := service.NewClientService(repo, cfg, kp)
	svc_contract := service.NewContractService(repo, cfg, kp)
	
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
    Core: hdl_core,
    User:   hdl_user,
    Client: hdl_client,
    Contract: hdl_contract,
}
	// Setup API routes
	router.SetupRoutes(r, allHandlers)

	logger.LogInfo("🎯 Server starting",
		zap.String("port", "8080"),
		zap.String("env", cfg.AppEnv),
		zap.String("version", Version),
		zap.String("commit", Commit),
		zap.String("buildTime", BuildTime),
	)

	// Start server
	if err := r.Run(":8080"); err != nil {
		logger.LogError("Server failed to start", err)
		log.Fatal(err)
	}
}

// printBanner prints startup banner
func printBanner() {
	fmt.Println("")
	fmt.Println("╔════════════════════════════════════════╗")
	fmt.Println("║     🚀 GOrders Backend Server 🚀     ║")
	fmt.Println("╠════════════════════════════════════════╣")
	fmt.Printf("║  Version:  %-30s ║\n", Version)
	fmt.Printf("║  Commit:   %-30s ║\n", Commit)
	fmt.Printf("║  Built:    %-30s ║\n", BuildTime)
	fmt.Println("╚════════════════════════════════════════╝")
	fmt.Println("")
}
