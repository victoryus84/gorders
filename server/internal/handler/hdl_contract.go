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

func (h *ContractHandler) CreateContract(c *gin.Context) {
	// 1. Parsarea corpului cererii (Folosim DTO-ul de Contract, NU de Client)
	// utils.ParseBody este o functie generica pe care ai definit-o in utils
	requests, err := utils.ParseBody[dto.ContractDTO](c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format JSON invalid: " + err.Error()})
		return
	}

	// 2. Luam ID-ul utilizatorului din token (pus de Middleware)
	val, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilizator neautorizat"})
		return
	}
	ownerID := val.(uint)

	// 3. LOGICA MASIVA S-A MUTAT IN SERVICE!
	// Handler-ul doar trimite datele si primeste rezultatul gata calculat
	result := h.svc.SyncContracts(requests, ownerID)

	// 4. Raspunsul catre 1C
	c.JSON(http.StatusCreated, result)
}

// 1. GetContractByID - Acum e METODĂ pe ContractHandler
func (h *ContractHandler) GetContractByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalid"})
		return
	}

	// Chemăm Serviciul (GetContractDetails - numele profesional pe care l-am ales)
	res, err := h.svc.GetContractDetails(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contractul nu a fost găsit"})
		return
	}

	c.JSON(http.StatusOK, res)
}

// 2. GetContractByClientID
func (h *ContractHandler) GetContractByClientID(c *gin.Context) {
	clientIDStr := c.Param("id")
	clientID, err := strconv.ParseUint(clientIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID client invalid"})
		return
	}

	// Serviciul livrează lista gata filtrată
	contracts, err := h.svc.GetClientContracts(uint(clientID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Eroare la recuperarea contractelor"})
		return
	}

	c.JSON(http.StatusOK, contracts)
}

// 3. CreateContractAddress (Importul de adrese)
func (h *ContractHandler) CreateContractAddress(c *gin.Context) {
	// Folosim utils.ParseBody cu DTO-ul specific
	requests, err := utils.ParseBody[dto.ContractAddressDTO](c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format invalid: " + err.Error()})
		return
	}

	val, _ := c.Get("user_id")
	ownerID := val.(uint)

	// TOATĂ logica de "for" și "create" s-a mutat în Service!
	// Handler-ul doar dă comanda de sincronizare.
	result := h.svc.RegisterContractAddress(requests, ownerID)

	c.JSON(http.StatusCreated, result)
}

// 4. GetContractAddressByID
func (h *ContractHandler) GetContractAddressByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID adresă invalid"})
		return
	}

	// Presupunem că am adăugat și această metodă în Service pentru detalii adresă
	res, err := h.svc.GetContractAddressDetails(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Adresa contractului nu a fost găsită"})
		return
	}

	c.JSON(http.StatusOK, res)
}