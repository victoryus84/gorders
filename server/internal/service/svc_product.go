package service

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/victoryus84/gorders/internal/config"
	"github.com/victoryus84/gorders/internal/dto"
	"github.com/victoryus84/gorders/internal/kafka"
	"github.com/victoryus84/gorders/internal/logger"
	"github.com/victoryus84/gorders/internal/models"
	"github.com/victoryus84/gorders/internal/repository"
)

// 1. INTERFAȚA PUBLICĂ
type ProductService interface {
	ProcessProductImport(requests []dto.ProductDTO) dto.ImportResult
	ProcessProductGroupImport(dtos []dto.ProductGroupDTO) dto.ImportResult
	GetAllProducts() ([]dto.ProductDTO, error)
}

// 2. STRUCTURA PRIVATĂ
type productService struct {
	rep_prd repository.ProductRepository
	rep_unt repository.UnitRepository
	rep_vat repository.VatTaxRepository
	cfg     *config.Config
	kfk     *kafka.Producer
}

// 3. CONSTRUCTORUL PENTRU UBER FX
func NewProductService(
	rep_prd repository.ProductRepository,
	rep_unt repository.UnitRepository,
	rep_vat repository.VatTaxRepository,
	cfg *config.Config,
	kfk *kafka.Producer) ProductService {
	return &productService{
		rep_prd: rep_prd,
		rep_unt: rep_unt,
		rep_vat: rep_vat,
		cfg:     cfg,
		kfk:     kfk,
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

func (svc *productService) ProcessProductGroupImport(requests []dto.ProductGroupDTO) dto.ImportResult {
	created := make([]*models.ProductGroup, 0)
	skipped := make([]map[string]string, 0)
	topic := svc.cfg.GetTopic("product_groups")

	// 1. INIȚIALIZAREA: Încărcăm grupele existente o singură dată
	groupMap := make(map[string]uint)
	if existingGroups, err := svc.rep_prd.GetAllProductGroups(); err == nil {
		for _, g := range existingGroups {
			groupMap[g.Code] = g.ID
		}
	}

	for _, req := range requests {
		// A. Validare de bază
		if strings.TrimSpace(req.Name) == "" {
			skipped = append(skipped, map[string]string{"name": req.Name, "reason": "missing_required_fields"})
			continue
		}

		// B. Verificare duplicate
		// existing, err := svc.rep.FindClientGroupByCode(req.Code)
		// if err == nil && existing != nil {
		//     skipped = append(skipped, map[string]string{"name": req.Name, "reason": "duplicate"})
		//     continue
		// }

		// C. MAPAREA IERARHIEI - Căutăm în dicționar dacă avem codul părintelui și luăm ID-ul lui
		var parentIDPtr *uint
		if cleanCode := strings.TrimSpace(req.ParentCode); cleanCode != "" && cleanCode != "not inserted" {
			if id, exists := groupMap[cleanCode]; exists {
				parentIDPtr = &id
			}
		}

		// D. Mapare DTO -> Model
		productgroup := &models.ProductGroup{
			Code:        req.Code,
			Name:        req.Name,
			Description: req.Description,
			ParentID:    parentIDPtr, // Acum ia adresa reală sau rămâne nil
		}

		// E. Salvare
		if err := svc.rep_prd.UpsertProductGroup(productgroup); err != nil {
			skipped = append(skipped, map[string]string{"code": req.Code, "reason": "upsert_failed: " + err.Error()})
			continue
		}

		// F. ACTUALIZAREA: Scriem grupa abia salvată în dicționar pentru viitorii ei copii!
		groupMap[productgroup.Code] = productgroup.ID

		// G. KAFKA
		go func(mod *models.ProductGroup) {
			payload, _ := json.Marshal(mod)
			_ = svc.kfk.Publish(context.Background(), topic, mod.Name, payload)
		}(productgroup)

		created = append(created, productgroup)
	}

	return dto.ImportResult{
		Status:         "success",
		TotalProcessed: len(created),
		TotalSkipped:   len(skipped),
		ErrorsPreview:  svc.limitErrors(skipped, 20),
		Message:        "Import finalizat",
	}
}

// Funcția care aduce toate produsele
func (svc *productService) GetAllProducts() ([]dto.ProductDTO, error) {
	products, err := svc.rep_prd.GetAllProducts()
	if err != nil {
		return nil, err
	}

	var response []dto.ProductDTO
	for _, p := range products {
		// 1. Extragem codul grupei în siguranță
		groupCode := ""
		if p.ProductGroupID != nil {
			groupCode = p.ProductGroup.Code
			// NOTĂ: Dacă Flutter are nevoie de NUMELE grupei, pui p.ProductGroup.Name
		}

		// 2. Protejăm pointer-ul de float64 (TVA) pentru a nu avea același pointer în toată lista
		vatValue := p.VatTax.Rate

		response = append(response, dto.ProductDTO{
			Code:        p.Code,
			Name:        p.Name,
			Description: p.Description,
			Article:     p.Article,
			Unit:        p.Unit.Name,
			VatCode:     p.VatTax.Code,
			VatTax:      &vatValue,
			GroupCode:   groupCode, // L-am legat aici
		})
	}
	return response, nil
}

// Funcția pentru Grupe de Produse

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
