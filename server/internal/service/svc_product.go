package service

import (
	"strings"

	"github.com/victoryus84/gorders/internal/dto"
	"github.com/victoryus84/gorders/internal/logger"
	"github.com/victoryus84/gorders/internal/models"
	"github.com/victoryus84/gorders/internal/repository"
)

// 1. INTERFAȚA PUBLICĂ
type ProductService interface {
	ProcessProductImport(requests []dto.ProductDTO) dto.ImportResult
	ProcessProductGroupImport(dtos []dto.ProductGroupDTO) map[string]interface{}
	GetFirst1000Products() ([]models.Product, error)
}

// 2. STRUCTURA PRIVATĂ
type productService struct {
	rep_prd repository.ProductRepository
	rep_unt repository.UnitRepository
	rep_vat repository.VatTaxRepository
}

// 3. CONSTRUCTORUL PENTRU UBER FX
func NewProductService(
	rep_prd repository.ProductRepository,
	rep_unt repository.UnitRepository,
	rep_vat repository.VatTaxRepository,
) ProductService {
	return &productService{
		rep_prd: rep_prd,
		rep_unt: rep_unt,
		rep_vat: rep_vat,
	}
}

// 4. IMPLEMENTAREA LOGICII DE SINCRONIZARE
func (svc *productService) ProcessProductImport(requests []dto.ProductDTO) dto.ImportResult {
	if len(requests) == 0 {
		return dto.ImportResult{
			Status:  "success",
			Message: "Niciun produs primit pentru procesare.",
		}
	}

	productsToSave := make([]models.Product, 0, len(requests))
	skipped := make([]map[string]string, 0)

	// Extragem dicționarele, cu tratare minimală a erorilor
	unitMap, err := svc.rep_unt.GetUnitIDMap()
	if err != nil {
		logger.LogError("Atenție: Eroare la încărcarea dicționarului de Unități", err)
	}

	vatMap, err := svc.rep_vat.GetVatIDMap()
	if err != nil {
		logger.LogError("Atenție: Eroare la încărcarea dicționarului de TVA", err)
	}

	groupMap, err := svc.rep_prd.GetGroupIDMap()
	if err != nil {
		logger.LogError("Atenție: Eroare la încărcarea dicționarului de Grupe Produse", err)
	}

	// Procesăm fiecare produs
	for _, req := range requests {
		// Validare
		if strings.TrimSpace(req.Code) == "" || strings.TrimSpace(req.Name) == "" {
			skipped = append(skipped, map[string]string{"code": req.Code, "reason": "missing_required_fields"})
			continue
		}

		// Mapare via Helpers
		dbUnitID := svc.resolveUnitID(req.Unit, unitMap)
		dbVatID := svc.resolveVatID(req.VatCode, vatMap)
		dbGroupID := svc.resolveGroupID(req.GroupCode, groupMap)

		// Construire Model
		product := models.Product{
			Code:           req.Code,
			Name:           req.Name,
			Description:    req.Description,
			Article:        req.Article,
			UnitID:         dbUnitID,
			VatTaxID:       dbVatID,
			ProductGroupID: dbGroupID,
		}

		productsToSave = append(productsToSave, product)
	}

	// Salvăm în masă
	if err := svc.rep_prd.UpsertProductsBatch(productsToSave, 500); err != nil {
		return dto.ImportResult{
			Status:  "error",
			Message: "Eroare fatală la salvarea în masă a produselor: " + err.Error(),
		}
	}

	return dto.ImportResult{
		Status:         "success",
		TotalProcessed: len(productsToSave),
		TotalSkipped:   len(skipped),
		ErrorsPreview:  svc.limitErrors(skipped, 20),
		Message:        "Sincronizare produse finalizată!",
	}
}

// Funcția care aduce primii 1000 de produse
func (svc *productService) GetFirst1000Products() ([]models.Product, error) {
	return svc.rep_prd.GetFirst1000Products()
}

// Funcția pentru Grupe de Produse
func (svc *productService) ProcessProductGroupImport(dtos []dto.ProductGroupDTO) map[string]interface{} {
	return map[string]interface{}{
		"status": "funcția pentru grupe este în construcție",
	}
}

// --- HELPER METODE PRIVATE ---

func (svc *productService) resolveUnitID(code string, unitMap map[string]uint) uint {
	if unitMap != nil {
		if id, ok := unitMap[strings.TrimSpace(code)]; ok && id != 0 {
			return id
		}
	}
	return 1
}

func (svc *productService) resolveVatID(code string, vatMap map[string]uint) uint {
	if vatMap != nil {
		if id, ok := vatMap[strings.TrimSpace(code)]; ok && id != 0 {
			return id
		}
	}
	return 1
}

func (svc *productService) resolveGroupID(code string, groupMap map[string]uint) *uint {
	if groupMap != nil {
		if id, ok := groupMap[strings.TrimSpace(code)]; ok && id != 0 {
			return &id
		}
	}
	return nil
}

func (svc *productService) limitErrors(skipped []map[string]string, limit int) []map[string]string {
	if len(skipped) > limit {
		return skipped[:limit]
	}
	return skipped
}