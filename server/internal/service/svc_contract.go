package service

import (
	"strings"
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
}

// ContractService - Ce oferim Handler-ului (GOrders API)
type ContractService interface {
	ProcessContractImport(requests []dto.ContractDTO, ownerID uint) dto.ImportResult
	GetContractsByClient(clientID uint) ([]dto.ContractDTO, error)
}

type contractService struct {
	rep ContractRepository
	clientRep ClientRepository
}

func NewContractService(rep ContractRepository, clientRep ClientRepository, cfg *config.Config, kp *kafka.Producer) ContractService {
	return &contractService{rep: rep, clientRep: clientRep}
}

// --- LOGICA DE BUSINESS (SYNC) ---

func (svc *contractService) ProcessContractImport(requests []dto.ContractDTO, ownerID uint) dto.ImportResult {
	skipped := make([]map[string]string, 0)

    // ==========================================
    // 1. DICȚIONARUL ÎN RAM (1 singur drum la DB!)
    // ==========================================
    clientMap, err := svc.clientRep.FindAllClientCodesMap()
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
            skipped = append(skipped, logSkip(req.Name, "Code (cod client) lipsă"))
            continue
        }

        dbClientID, exists := clientMap[req.Code]
        if !exists {
            skipped = append(skipped, logSkip(req.Name, "Client inexistent: "+req.Code))
            continue
        }

     // Conversie DTO -> Model (asigură-te că în mapDTOToModel asignezi contract.SyncID = req.SyncID)
        contract := svc.mapDTOToModel(req, dbClientID, ownerID)
        
       // Get the key directly from 1C
        syncKey := req.SyncID
        
        // Check if 1C sent a duplicate in the same file
        if _, exists := uniqueContracts[syncKey]; exists {
            // Log in English for the server console
            fmt.Printf("🔥 XML DUPLICATE: Contract with SyncID '%s' appeared multiple times in the payload!\n", syncKey)
            
            // Textul pentru 1C rămâne în română (fără diacritice) ca să-l înțeleagă operatorii
            skipped = append(skipped, logSkip(req.Name, "Dublura in XML: "+syncKey))
        }
        
        // Add to the map
        uniqueContracts[syncKey] = contract
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
		SyncID: req.SyncID,
		Name:     req.Name,
		Number:   stringPtr(req.Number),
		Date:     parse1CDate(req.Date),
		Amount:   req.Amount,
		Status:   req.Status,
		ClientID: clientID,
		OwnerID:  ownerID,
	}
}

func (svc *contractService) limitErrors(skipped []map[string]string, limit int) []map[string]string {
	if len(skipped) > limit {
		return skipped[:limit]
	}
	return skipped
}
