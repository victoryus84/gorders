package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handlers - "Tava" universală care grupează toți chelnerii (specialiștii)
type Handlers struct {
	Core   *CoreHandler
	User   *UserHandler
	Client *ClientHandler
	// Aici vei adăuga Order, Product, etc. pe viitor
}

// CoreHandler - Specialistul pentru starea sistemului (fostul Health)
type CoreHandler struct {
	db      *gorm.DB
	version string
	commit  string
}

// NewCoreHandler - Constructorul pentru Core
func NewCoreHandler(db *gorm.DB, version, commit string) *CoreHandler {
	return &CoreHandler{
		db:      db,
		version: version,
		commit:  commit,
	}
}

// Check - Endpoint pentru sănătatea sistemului (/health)
func (h *CoreHandler) Check(c *gin.Context) {
	status := "healthy"
	dbStatus := "ok"

	// Verificăm conexiunea la baza de date
	sqlDB, err := h.db.DB()
	if err != nil || sqlDB.Ping() != nil {
		status = "unhealthy"
		dbStatus = "connection failed"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    status,
		"database":  dbStatus,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"version":   h.version,
	})
}

// Version - Endpoint pentru versiunea aplicației (/version)
func (h *CoreHandler) Version(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version": h.version,
		"commit":  h.commit,
	})
}

// Ping - Un simplu test de latență (/ping)
func (h *CoreHandler) Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}