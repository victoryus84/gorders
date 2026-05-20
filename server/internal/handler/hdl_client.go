package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/victoryus84/gorders/internal/dto"
	"github.com/victoryus84/gorders/internal/service"
	"github.com/victoryus84/gorders/internal/utils"
)

type ClientHandler struct {
	svc service.ClientService
}

func NewClientHandler(svc service.ClientService) *ClientHandler {
	return &ClientHandler{svc: svc}
}

// CreateClient - Handler pentru import masiv de clienți (1C style)
func (h *ClientHandler) CreateClient(ctx *gin.Context) {
	// 1. Parsăm body-ul (JSON/XML/Array) folosind utilitarul tău deștept
	requests, err := utils.ParseBody[dto.ClientDTO](ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Format invalid: " + err.Error()})
		return
	}

	// 2. Trimitem tot calupul la Service (Bucătarul se ocupă de validări și duplicate)
	result := h.svc.ProcessClientImport(requests)

	// 3. Răspunsul final
	ctx.JSON(http.StatusCreated, result)
}

func (hdl *ClientHandler) CreateClientGroup(ctx *gin.Context) {
	// 1. Parsăm body-ul (JSON/XML/Array) folosind utilitarul tău deștept
	requests, err := utils.ParseBody[dto.ClientGroupDTO](ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Format invalid: " + err.Error()})
		return
	}

	// 2. Trimitem tot calupul la Service (Bucătarul se ocupă de validări și duplicate)
	result := hdl.svc.ProcessClientGroupImport(requests)

	// 3. Răspunsul final
	ctx.JSON(http.StatusCreated, result)
}

// GetClients - Obține primii 1000 de clienți
func (hdl *ClientHandler) GetClients(ctx *gin.Context) {
	clients, err := hdl.svc.SvcGetFirst1000Clients()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, clients)
}

// SearchClients - Căutare clienți după query (q=...)
func (hdl *ClientHandler) SearchClients(ctx *gin.Context) {
	query := ctx.Query("q")
	if query == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "The parameter 'q' is required"})
		return
	}

	// Serviciul ne întoarce direct DTO-urile pregătite pentru export (cu email string, nu pointer)
	response, err := hdl.svc.SearchClients(query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetClientByID - Obținerea unui singur client
func (hdl *ClientHandler) GetClientByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	client, err := hdl.svc.SearchClientByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Clientul nu a fost găsit"})
		return
	}

	ctx.JSON(http.StatusOK, client)
}

// CreateClientAddress - Import masiv de adrese pentru clienți
func (hdl *ClientHandler) CreateClientAddress(c *gin.Context) {
	requests, err := utils.ParseBody[dto.ClientAddressDTO](c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format invalid"})
		return
	}

	// Luăm ID-ul utilizatorului din token (pus acolo de Middleware-ul de Auth)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilizator neidentificat"})
		return
	}

	// Trimitem la service adresele și cine e proprietarul lor
	result := hdl.svc.ProcessAddressImport(requests, userID.(uint))

	c.JSON(http.StatusCreated, result)
}

