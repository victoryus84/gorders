package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/victoryus84/gorders/internal/dto"
	"github.com/victoryus84/gorders/internal/service"

)

type ContractHandler struct {
	svc service.ContractService
}	

func NewContractHandler(svc service.ContractService) *ContractHandler {
	return &ContractHandler{svc: svc}
}

func CreateContract(h *ContractHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		// PASUL 1: Parsarea corpului cererii (JSON sau XML)
		// Folosim DTO-ul specific pentru contracte
		requests, err := ParseBody[ContractReq](c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format invalid: " + err.Error()})
			return
		}

		// Luăm ID-ul utilizatorului din token (pentru coloana owner_id din DB)
		val, _ := c.Get("user_id")
		ownerID := val.(uint)

		created := make([]*models.Contract, 0)
		skipped := make([]map[string]string, 0)

		// PASUL 2: Logica de business pentru fiecare contract
		for _, req := range requests {
			
			// A. Validare detaliată
			// Verificăm dacă lipsește Codul Fiscal
			if strings.TrimSpace(req.FiscalID) == "" {
				skipped = append(skipped, map[string]string{
					"number": req.Name,
					"reason": "EROARE: Campul 'fiscal_id' a ajuns GOL. Clientul nu poate fi gasit!",
				})
				continue
			}

			// B. Căutarea clientului (rămâne la fel)
			client, err := s.FindClientByFiscalID(req.FiscalID)
			if err != nil {
				skipped = append(skipped, map[string]string{
					"number": req.Name,
					"reason": "Clientul cu FiscalID " + req.FiscalID + " nu exista in baza de date!",
				})
				continue
			}

			// --- LOGICA PENTRU NUMBER (NUMAR SAU NULL) ---
			rawNumber := strings.TrimSpace(req.Number)
			var numberPtr *string
			// Dacă e număr pe bune, salvăm numărul
			if rawNumber != "" {
				numberPtr = &rawNumber
			}

			// --- LOGICA PENTRU DATA (INFINITĂ SAU NORMALĂ) ---
			var datePtr *time.Time // Default este nil (NULL în Postgres)
			dateStr := strings.TrimSpace(req.Date)

			if dateStr != "" && dateStr != "00.00.0000" && dateStr != "01-01-0001" {
				// Parsăm data folosind formatul tău din 1C (Zi-Lună-An)
				t, err := time.Parse("02-01-2006", dateStr)
				if err == nil {
					datePtr = &t
				} else {
					// Folosim log.Printf pentru a vedea eroarea în consola Gin
					log.Printf("[IMPORT ERROR] Data invalida pentru contractul %s: %v", req.Number, err)
				}
			}

			// C. Conversia de la REQ (ce vine din 1C) la MODEL (ce pleacă în Postgres)
			contract := &models.Contract{
				Name:     req.Name,
				Number:   numberPtr,
				Date:     datePtr,
				Amount:   req.Amount,
				Status:   req.Status,
				ClientID: client.ID, // <--- Aici e "magia": ID-ul de Postgres al clientului
				OwnerID:  ownerID,   // Coloana owner_id pe care o aveai în DB
			}

			// D. Salvarea efectivă prin Service
			if err := s.CreateContract(contract); err != nil {
				skipped = append(skipped, map[string]string{
					"number": req.Number,
					"reason": "db_save_error: " + err.Error(),
				})
				continue
			}
			created = append(created, contract)
		}

		// --- PASUL 3: Răspuns "LIGHT" pentru 1C 7.7 ---

		// Pregătim o listă scurtă cu erorile (doar primele 10, să nu omorâm 1C-ul)
		shortSkipped := skipped
		if len(skipped) > 20 {
			shortSkipped = skipped[:20]
		}

		c.JSON(http.StatusCreated, gin.H{
			"status":         "success",
			"total_created":  len(created),
			"total_skipped":  len(skipped),
			"errors_preview": shortSkipped, // Trimitem doar o mostră de erori
			"message":        "Import contracts finalizat cu succes",
		})
	}
}

func GetContractByID(h *ContractHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		res, err := s.FindContractByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Contract not found"})
			return
		}
		c.JSON(http.StatusOK, res)
	}
}

func GetContractByClientID(s Service) gin.HandlerFunc {
return func(c *gin.Context) {
        // 1. Extrage ID-ul clientului din URL
        clientIDStr := c.Param("id")
        
        // 2. Convertește-l la tipul tău de date (de obicei uint sau int)
        clientID, err := strconv.ParseUint(clientIDStr, 10, 64)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "ID client invalid"})
            return
        }

        // 3. Apeleză serviciul pentru a aduce datele
        // Presupunem că metoda se numește GetByClientID și returnează []Contract, error
        contracts, err := s.FindContractByClientID(uint(clientID))
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Nu s-au putut recupera contractele"})
            return
        }

        // 4. Returnează lista de contracte
        c.JSON(http.StatusOK, contracts)
    }
}
// Contract Address handlers

func CreateContractAddress(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		requests, err := ParseBody[ContractAddressReq](c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format"})
			return
		}

		val, _ := c.Get("user_id")
		ownerID := val.(uint)

		created := make([]*models.ContractAddress, 0)
		for _, req := range requests {
			addr := &models.ContractAddress{
				ContractID: req.ContractID,
				Address:    req.Address,
				Type:       req.Type,
				OwnerID:    ownerID,
			}
			if err := s.CreateContractAddress(addr); err == nil {
				created = append(created, addr)
			}
		}
		c.JSON(http.StatusCreated, created)
	}
}

func GetContractAddressByID(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		res, err := s.FindContractAddressByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
			return
		}
		c.JSON(http.StatusOK, res)
	}
}
