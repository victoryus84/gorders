package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/victoryus84/gorders/internal/dto"
	"github.com/victoryus84/gorders/internal/service"
	"github.com/victoryus84/gorders/internal/utils"
)

type ContractHandler struct {
	svc service.ContractService
}

func NewContractHandler(svc service.ContractService) *ContractHandler {
	return &ContractHandler{svc: svc}
}

// POST /contracts
func (h *ContractHandler) CreateContract(c *gin.Context) {
	// Folosim ParseBody pentru a lua DTO-ul (asigură-te că în dto.ContractDTO ai tag-ul binding:"required" pe ClientID)
	request, err := utils.ParseBody[dto.ContractDTO](c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date invalide: " + err.Error()})
		return
	}

	ownerID, ok := h.getUserID(c)
	if !ok { return }

	// Notă: SyncContracts primește de obicei un slice. 
	// Dacă trimiți un singur obiect, verifică dacă serviciul așteaptă []dto.ContractDTO sau doar unul.
	result := h.svc.SyncContracts(request, ownerID)
	c.JSON(http.StatusCreated, result)
}

// GET /contracts/:client_id
func (h *ContractHandler) GetContractsByClientID(c *gin.Context) {
	// Extragem "client_id" direct din URL-ul definit în rute
	clientID, ok := h.parseID(c, "client_id")
	if !ok { return }

	contracts, err := h.svc.GetContractsByClient(clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Eroare la recuperarea listei de contracte"})
		return
	}

	c.JSON(http.StatusOK, contracts)
}

// --- HELPER FUNCTIONS ---

func (h *ContractHandler) getUserID(c *gin.Context) (uint, bool) {
	val, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Lipsă autorizare"})
		return 0, false
	}
	
	ownerID, ok := val.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tip date user_id invalid"})
		return 0, false
	}
	
	return ownerID, true
}

func (h *ContractHandler) parseID(c *gin.Context, paramName string) (uint, bool) {
	idStr := c.Param(paramName)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalid: trebuie să fie un număr pozitiv"})
		return 0, false
	}
	return uint(id), true
}