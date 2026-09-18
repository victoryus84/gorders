package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/victoryus84/gorders/internal/config"
)

// Handlers - "Tava" universală care grupează toți chelnerii (specialiștii)
type Handlers struct {
	Core   *CoreHandler
	User   *UserHandler
	Client *ClientHandler
	Contract *ContractHandler
	Product *ProductHandler
	// Aici vei adăuga Order, Product, etc. pe viitor
}

// CoreHandler - Specialistul pentru starea sistemului (fostul Health)
type CoreHandler struct {
    db  *gorm.DB
    cfg *config.Config
}

// NewCoreHandler - Constructorul pentru Core
func NewCoreHandler(db *gorm.DB, cfg *config.Config) *CoreHandler {
    return &CoreHandler{
        db:  db,
        cfg: cfg, // Acum ai acces la cfg.Version și cfg.Commit în rutele tale (ex: la /health)
    }
}

// Check - Endpoint pentru sănătatea sistemului (/health)
func (hdl *CoreHandler) Check(ctx *gin.Context) {
	status := "healthy"
	dbStatus := "ok"

	// Verificăm conexiunea la baza de date
	sqlDB, err := hdl.db.DB()
	if err != nil || sqlDB.Ping() != nil {
		status = "unhealthy"
		dbStatus = "connection failed"
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":    status,
		"database":  dbStatus,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"version":   hdl.cfg.Version,
	})
}

// Version - Endpoint pentru versiunea aplicației (/version)
func (hdl *CoreHandler) Version(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"version": hdl.cfg.Version,
		"commit":  hdl.cfg.Commit,
	})
}

// Ping - Un simplu test de latență (/ping)
func (hdl *CoreHandler) Ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, "pong")
}

func RegisterCoreRoutes(rte *gin.Engine, hdl *CoreHandler) {
	rte.GET("/health", hdl.Check)
	rte.GET("/version", hdl.Version)
	rte.GET("/ping", hdl.Ping)
}