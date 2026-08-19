package service

import (
	"strings"
	"time"
	"fmt"
	"github.com/victoryus84/gorders/internal/config"
	"github.com/victoryus84/gorders/internal/dto"
	"github.com/victoryus84/gorders/internal/kafka"
	"github.com/victoryus84/gorders/internal/models"
)

// ContractRepository - Ce așteptăm de la baza de date
type ContractRepository interface {
	CreateContract(contract *models.Contract) error
	UpsertContractBatch(contracts []*models.Contract, batchSize int) error
	FindContractsByClientID(clientID uint) ([]models.Contract, error)
	FindClientByCode(code string) (*models.Client, error)
	FindContractByID(id uint) (*models.Contract, error)
	FindAllClientCodesMap() (map[string]uint, error)
}

// ContractService - Ce oferim Handler-ului (GOrders API)
type ContractService interface {
	ProcessContractImport(requests []dto.ContractDTO, ownerID uint) dto.ImportResult
	GetContractDetails(id uint) (*dto.ContractDTO, error) // Am aliniat numele cu Handler-ul
	GetContractsByClient(clientID uint) ([]dto.ContractDTO, error)
}

type contractService struct {
	rep ContractRepository
}

func NewContractService(rep ContractRepository, cfg *config.Config, kp *kafka.Producer) ContractService {
	return &contractService{rep: rep}
}

// --- LOGICA DE BUSINESS (SYNC) ---

func (svc *contractService) ProcessContractImport(requests []dto.ContractDTO, ownerID uint) dto.ImportResult {
	skipped := make([]map[string]string, 0)

    // ==========================================
    // 1. DICȚIONARUL ÎN RAM (1 singur drum la DB!)
    // ==========================================
    clientMap, err := svc.rep.FindAllClientCodesMap()
    if err != nil {
        clientMap = make(map[string]uint) 
    }

    // AICI E "SITA": Pregătim un map în loc de array pentru a elimina dublurile automat
    uniqueContracts := make(map[string]*models.Contract)

    // ==========================================
    // 2. PROCESĂM ÎN MEMORIE ȘI FILTRĂM DUBURILE
    // ==========================================
    for _, req := range requests {

        if strings.TrimSpace(req.Code) == "" {
            skipped = append(skipped, svc.logSkip(req.Name, "Code (cod client) lipsă"))
            continue
        }

        dbClientID, exists := clientMap[req.Code]
        if !exists {
            skipped = append(skipped, svc.logSkip(req.Name, "Client inexistent: "+req.Code))
            continue
        }

      // Conversie DTO -> Model
        contract := svc.mapDTOToModel(req, dbClientID, ownerID)
        
        // --- GENERĂM IDENTIFICATORUL EXTERN UNIC ---
        // Ex: "00015_00002"
        contract.SyncID = fmt.Sprintf("%s_%s", req.Code, req.Number)

        // Sita din RAM va folosi tot această cheie ca să prevină dublurile din același XML!
        uniqueContracts[contract.SyncID] = contract
    }

    // --- RECONSTRUIM ARRAY-UL PENTRU GORM ---
    // Acum că am scăpat de dubluri, mutăm datele din map înapoi într-un slice
    contractsToSave := make([]*models.Contract, 0, len(uniqueContracts))
    for _, c := range uniqueContracts {
        contractsToSave = append(contractsToSave, c)
    }

    // ==========================================
    // 3. SALVAREA ÎN MASĂ (Tunul!)
    // ==========================================
    if len(contractsToSave) > 0 {
        if err := svc.rep.UpsertContractBatch(contractsToSave, 1000); err != nil {
            return dto.ImportResult{
                Status:  "error",
                Message: "Eroare DB la salvarea în masă: " + err.Error(),
            }
        }
    }

    return dto.ImportResult{
        Status:         "success",
        TotalProcessed: len(contractsToSave), // Acum arată nr real, fără dubluri!
        TotalSkipped:   len(skipped),
        ErrorsPreview:  svc.limitErrors(skipped, 20),
        Message:        "Sincronizare contracte finalizată instantaneu!",
    }
}

// --- IMPLEMENTARE METODE LIPSĂ (Pentru a repara erorile de compilare) ---

func (svc *contractService) GetContractDetails(id uint) (*dto.ContractDTO, error) {
	contract, err := svc.rep.FindContractByID(id)
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

func (svc *contractService) GetContractsByClient(clientID uint) ([]dto.ContractDTO, error) {
	contracts, err := svc.rep.FindContractsByClientID(clientID)
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

func (svc *contractService) mapDTOToModel(req dto.ContractDTO, clientID uint, ownerID uint) *models.Contract {
	return &models.Contract{
		Name:     req.Name,
		Number:   svc.stringPtr(req.Number),
		Date:     svc.parse1CDate(req.Date),
		Amount:   req.Amount,
		Status:   req.Status,
		ClientID: clientID,
		OwnerID:  ownerID,
	}
}

func (svc *contractService) parse1CDate(dateStr string) *time.Time {
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

func (svc *contractService) stringPtr(str string) *string {
	str = strings.TrimSpace(str)
	if str == "" {
		return nil
	}
	return &str
}

func (svc *contractService) logSkip(name, reason string) map[string]string {
	return map[string]string{"item": name, "reason": reason}
}

func (svc *contractService) limitErrors(skipped []map[string]string, limit int) []map[string]string {
	if len(skipped) > limit {
		return skipped[:limit]
	}
	return skipped
}
