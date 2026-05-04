package service

import (
	"strings"
	"time"

	"github.com/victoryus84/gorders/internal/dto"
	"github.com/victoryus84/gorders/internal/models"
)

// ContractRepository - Ce așteptăm de la baza de date
type ContractRepository interface {
    CreateContract(contract *models.Contract) error
    FindContractsByClientID(clientID uint) ([]models.Contract, error)
    FindClientByFiscalID(fiscalID string) (*models.Client, error)
	FindContractByID(id uint) (*models.Contract, error)
}
// ContractService - Ce oferim Handler-ului (GOrders API)
type ContractService interface {
	SyncContracts(requests []dto.ContractDTO, ownerID uint) dto.ImportResult
	GetContractDetails(id uint) (*dto.ContractDTO, error) // Am aliniat numele cu Handler-ul
	GetContractsByClient(clientID uint) ([]dto.ContractDTO, error)
}

type contractService struct {
	repo ContractRepository
}

func NewContractService(repo ContractRepository) ContractService {
	return &contractService{repo: repo}
}

// --- LOGICA DE BUSINESS (SYNC) ---

func (s *contractService) SyncContracts(requests []dto.ContractDTO, ownerID uint) dto.ImportResult {
	createdCount := 0
	skipped := make([]map[string]string, 0)

	for _, req := range requests {
		// 1. Validare rapidă
		if strings.TrimSpace(req.FiscalID) == "" {
			skipped = append(skipped, s.logSkip(req.Name, "FiscalID lipsă"))
			continue
		}

		// 2. Găsire client
		client, err := s.repo.FindClientByFiscalID(req.FiscalID)
		if err != nil {
			skipped = append(skipped, s.logSkip(req.Name, "Client inexistent: "+req.FiscalID))
			continue
		}

		// 3. Conversie DTO -> Model (Logica de transformare e izolată)
		contract := s.mapDTOToModel(req, client.ID, ownerID)

		// 4. Salvare
		if err := s.repo.CreateContract(contract); err != nil {
			skipped = append(skipped, s.logSkip(req.Number, "Eroare DB: "+err.Error()))
			continue
		}
		createdCount++
	}

	return dto.ImportResult{
		Status:        "success",
		TotalCreated:  createdCount,
		TotalSkipped:  len(skipped),
		ErrorsPreview: s.limitErrors(skipped, 20),
		Message:       "Sincronizare contracte finalizată",
	}
}

// --- IMPLEMENTARE METODE LIPSĂ (Pentru a repara erorile de compilare) ---

func (s *contractService) GetContractDetails(id uint) (*dto.ContractDTO, error) {
	contract, err := s.repo.FindContractByID(id)
	if err != nil {
		return nil, err
	}
	// Aici transformi modelul de DB înapoi în DTO pentru JSON
	return &dto.ContractDTO{
		Name:   contract.Name,
		Amount: contract.Amount,
		Status: contract.Status,
		// Adaugă restul câmpurilor necesare
	}, nil
}

func (s *contractService) GetContractsByClient(clientID uint) ([]dto.ContractDTO, error) {
	contracts, err := s.repo.FindContractsByClientID(clientID)
	if err != nil {
		return nil, err
	}

	var result []dto.ContractDTO
	for _, c := range contracts {
		result = append(result, dto.ContractDTO{
			Name:   c.Name,
			Amount: c.Amount,
		})
	}
	return result, nil
}

// --- HELPER FUNCTIONS (Curățenia din labirint) ---

func (s *contractService) mapDTOToModel(req dto.ContractDTO, clientID uint, ownerID uint) *models.Contract {
	return &models.Contract{
		Name:     req.Name,
		Number:   s.stringPtr(req.Number),
		Date:     s.parse1CDate(req.Date),
		Amount:   req.Amount,
		Status:   req.Status,
		ClientID: clientID,
		OwnerID:  ownerID,
	}
}

func (s *contractService) parse1CDate(dateStr string) *time.Time {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" || dateStr == "00.00.0000" || dateStr == "01-01-0001" {
		return nil
	}
	t, err := time.Parse("02-01-2006", dateStr)
	if err != nil {
		return nil
	}
	return &t
}

func (s *contractService) stringPtr(str string) *string {
	str = strings.TrimSpace(str)
	if str == "" {
		return nil
	}
	return &str
}

func (s *contractService) logSkip(name, reason string) map[string]string {
	return map[string]string{"item": name, "reason": reason}
}

func (s *contractService) limitErrors(skipped []map[string]string, limit int) []map[string]string {
	if len(skipped) > limit {
		return skipped[:limit]
	}
	return skipped
}
