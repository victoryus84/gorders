package service

import (
	"strings"
	"time"

	"github.com/victoryus84/gorders/internal/dto"
	"github.com/victoryus84/gorders/internal/models"
)

type ContractRepository interface {
	FindClientByFiscalID(fiscalID string) (*models.Client, error)
	CreateContract(contract *models.Contract) error
	FindContractByID(id uint) (*models.Contract, error)
	FindContractsByClientID(clientID uint) ([]models.Contract, error)
	CreateContractAddress(address *models.ContractAddress) error
	FindContractAddressByID(id uint) (*models.ContractAddress, error)
}

type ContractService interface {
	SyncContracts(requests []dto.ContractDTO, ownerID uint) dto.ImportResult
	FindContractByID(id uint) (*dto.ContractDTO, error)
	FindContractByClientID(clientID uint) ([]dto.ContractDTO, error)
	CreateContractAddress(requests []dto.ContractAddressDTO, ownerID uint) dto.ImportResult
	FindContractAddressByID(id uint) (*dto.ContractAddressDTO, error)
}

type contractService struct {
	repo ContractRepository
}		

func NewContractService(repo ContractRepository) ContractService {
	return &contractService{repo: repo}
}	

func (s *contractService) SyncContracts(requests []dto.ContractDTO, ownerID uint) dto.ImportResult {
	createdCount := 0
	skipped := make([]map[string]string, 0)

	for _, req := range requests {
		// A. Validare: Codul Fiscal este obligatoriu pentru a găsi clientul
		if strings.TrimSpace(req.FiscalID) == "" {
			skipped = append(skipped, map[string]string{
				"number": req.Name,
				"reason": "EROARE: FiscalID lipsa. Nu pot asocia contractul!",
			})
			continue
		}

		// B. Căutarea clientului în DB prin Repository
		client, err := s.repo.FindClientByFiscalID(req.FiscalID)
		if err != nil {
			skipped = append(skipped, map[string]string{
				"number": req.Name,
				"reason": "Clientul cu FiscalID " + req.FiscalID + " nu exista!",
			})
			continue
		}

		// C. Logica pentru Number (Păstrăm pointerul pentru NULL în DB)
		rawNumber := strings.TrimSpace(req.Number)
		var numberPtr *string
		if rawNumber != "" {
			numberPtr = &rawNumber
		}

		// D. Logica pentru Data (Sanitizare format 1C)
		var datePtr *time.Time
		dateStr := strings.TrimSpace(req.Date)
		if dateStr != "" && dateStr != "00.00.0000" && dateStr != "01-01-0001" {
			t, err := time.Parse("02-01-2006", dateStr)
			if err == nil {
				datePtr = &t
			}
		}

		// E. Conversie DTO -> Model
		contract := &models.Contract{
			Name:     req.Name,
			Number:   numberPtr,
			Date:     datePtr,
			Amount:   req.Amount,
			Status:   req.Status,
			ClientID: client.ID, // ID-ul real din Postgres
			OwnerID:  ownerID,
		}

		// F. Salvarea prin Repository
		if err := s.repo.CreateContract(contract); err != nil {
			skipped = append(skipped, map[string]string{
				"number": req.Number,
				"reason": "db_error: " + err.Error(),
			})
			continue
		}
		createdCount++
	}

	return dto.ImportResult{
		Status:        "success",
		TotalCreated:  createdCount,
		TotalSkipped:  len(skipped),
		ErrorsPreview: s.limitErrors(skipped, 20),
		Message:       "Import contracte finalizat",
	}
}

// Helper pentru limitarea erorilor (pune-l la finalul fișierului)
func (s *contractService) limitErrors(skipped []map[string]string, limit int) []map[string]string {
	if len(skipped) > limit {
		return skipped[:limit]
	}
	return skipped
}