package service

import (
	"strings"

	"github.com/victoryus84/gorders/internal/dto"
	"github.com/victoryus84/gorders/internal/models"
	"github.com/victoryus84/gorders/internal/repository"
)

// 1. INTERFAȚA PUBLICĂ
type ValuteService interface {
	ProcessImport(requests []dto.ValuteDTO) dto.ImportResult
	FindAll() ([]models.Valute, error)
}

// 2. STRUCTURA PRIVATĂ
type valuteService struct {
	repo repository.ValuteRepository
}

func NewValuteService(repo repository.ValuteRepository) ValuteService {
	return &valuteService{repo: repo}
}

// 3. IMPLEMENTAREA LOGICII PRINCIPALE
func (svc *valuteService) ProcessImport(requests []dto.ValuteDTO) dto.ImportResult {
	if len(requests) == 0 {
		return dto.ImportResult{Status: "success", Message: "Niciun element primit."}
	}

	itemsToSave := make([]models.Valute, 0, len(requests))
	skipped := make([]map[string]any, 0)


	for _, req := range requests {
		// Validare de bază
		if strings.TrimSpace(req.Code) == "" {
			skipped = append(skipped, map[string]any{"code": req.Code, "reason": "missing_code"})
			continue
		}

		item := models.Valute{
			Name: req.Name,
			Description: req.Description,
			Code: req.Code,
			Symbol: req.Symbol,
			IsBase: req.IsBase,
		}


		itemsToSave = append(itemsToSave, item)
	}

	if err := svc.repo.UpsertBatch(itemsToSave, 500); err != nil {
		return dto.ImportResult{Status: "error", Message: "Eroare la salvare: " + err.Error()}
	}

	return dto.ImportResult{
		Status:         "success",
		TotalProcessed: len(itemsToSave),
		TotalSkipped:   len(skipped),
		Message:        "Sincronizare Valute finalizată!",
	}
}

func (svc *valuteService) FindAll() ([]models.Valute, error) {
	return svc.repo.FindAll()
}


