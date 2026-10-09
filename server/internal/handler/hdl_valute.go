package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/victoryus84/gorders/internal/dto"
	"github.com/victoryus84/gorders/internal/logger"
	"github.com/victoryus84/gorders/internal/service"
	"github.com/victoryus84/gorders/internal/utils"
)

type ValuteHandler struct {
	svc service.ValuteService
}

func NewValuteHandler(svc service.ValuteService) *ValuteHandler {
	return &ValuteHandler{svc: svc}
}

// Înregistrarea rutelor pentru Uber Fx (module.go)
func RegisterValuteRoutes(r *gin.Engine, h *ValuteHandler) {
	api := r.Group("/api/valutes")
	{
		api.POST("/import", h.Create)
	}
}

// Sincronizare Entități Principale
func (h *ValuteHandler) Create(c *gin.Context) {
	requests, err := utils.ParseBody[dto.ValuteDTO](c)
	if err != nil {
		logger.LogError("⚠️ Format invalid primit de la 1C pentru Valute", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format invalid: " + err.Error()})
		return
	}

	result := h.svc.ProcessImport(requests)
	
	status := http.StatusCreated
	if result.Status == "error" {
		status = http.StatusInternalServerError
	}
	c.JSON(status, result)
}


