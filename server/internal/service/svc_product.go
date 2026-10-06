package service

import (
	"context"
	"encoding/json"
	"strings"
	"time"

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
	ProcessProductPriceImport(dtos []dto.ProductPriceDTO) dto.ImportResult
	GetProducts() ([]dto.ProductDTO, error)
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
	skipped := make([]map[string]any, 0)

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
			skipped = append(skipped, map[string]any{"code": req.Code, "reason": "missing_required_fields"})
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
	skipped := make([]map[string]any, 0)
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
			skipped = append(skipped, map[string]any{"name": req.Name, "reason": "missing_required_fields"})
			continue
		}

		// B. Verificare duplicate
		// existing, err := svc.rep.FindClientGroupByCode(req.Code)
		// if err == nil && existing != nil {
		//     skipped = append(skipped, map[string]any{"name": req.Name, "reason": "duplicate"})
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
			skipped = append(skipped, map[string]any{"code": req.Code, "reason": "upsert_failed: " + err.Error()})
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

func (svc *productService) ProcessProductPriceImport(requests []dto.ProductPriceDTO) dto.ImportResult {
    // Am corectat PriceProduct în ProductPrice ca să se potrivească
    created := make([]*models.PriceProduct, 0)
    skipped := make([]map[string]any, 0)
    
    // topic := svc.cfg.GetTopic("product_prices") // Decomentează când implementezi trimiterea către Kafka

   for _, req := range requests {
    // A. Validare de bază
    if req.Price == 0 {
        skipped = append(skipped, map[string]any{
			"price": req.ProductCode, 
			"reason": "missing price",
		})
        continue
    }

    // --- NOU: LOGICA DE CONVERSIE A DATEI ---
    var dateFromPtr *time.Time // Implicit este nil (null)
    if req.DateFrom != "" {
        // Transformăm textul în time.Time. 
        // time.RFC3339 așteaptă formatul standard gen: "2026-10-05T15:04:05Z"
        parsedTime, err := time.Parse(time.RFC3339, req.DateFrom)
        if err != nil {
            // Dacă 1C trimite o dată greșită, sărim peste produs
            skipped = append(skipped, map[string]any{
				"code": req.ProductCode, 
				"reason": "invalid_date_format",
			})
            continue
        }
        dateFromPtr = &parsedTime // Luăm adresa memoriei (pointerul) ca să se potrivească cu *time.Time
    }

    // C. Mapare DTO -> Model
    productprice := &models.PriceProduct{
        Price:    req.Price,
        DateFrom: dateFromPtr, // <--- Acum e pointerul corect! Dacă textul era gol, va fi nil.
    }
        // D. Salvare (unul câte unul) - Dacă ai multe, pe viitor facem aici acel UpsertBatch!
        if err := svc.rep_prd.UpsertProductPrice(productprice); err != nil {
            skipped = append(skipped, map[string]any{
				"code": req.ProductCode, 
				"reason": "upsert_failed: " + err.Error(),
			})
            continue
        }

        created = append(created, productprice)
    }

    // AICI se termină funcția de import!
    return dto.ImportResult{
        Status:         "success",
        TotalProcessed: len(created),
        TotalSkipped:   len(skipped),
        ErrorsPreview:  svc.limitErrors(skipped, 20),
        Message:        "Import prețuri finalizat",
    }
}

// GetProducts aduce tot nomenclatorul și îl mapează în DTO-uri pentru Flutter
func (svc *productService) GetProducts() ([]dto.ProductDTO, error) {
	// 1. Aducem modelele brute din DB (care includ relațiile via Preload)
	// Notă: Dacă în repozitoriul tău metoda se numește altfel (ex: GetAllProducts), modifică aici:
	products, err := svc.rep_prd.GetAllProducts() 
	if err != nil {
		return nil, err
	}

	// 2. Pregătim lista de răspuns
	var response []dto.ProductDTO

	// 3. Mapăm fiecare model GORM în ProductDTO
	for _, p := range products {
		// Protecție: Extragem codul grupei doar dacă produsul are o grupă atașată
		groupCode := ""
		if p.ProductGroupID != nil {
			groupCode = p.ProductGroup.Code 
		}

		// Protecție: Salvăm valoarea TVA într-o variabilă locală pentru a prelua adresa (&) corect
		// ATENȚIE: Schimbă '.Value' cu numele real din models.VatTax (ex: .Tax, .Rate)
		vatValue := p.VatTax.Rate

		response = append(response, dto.ProductDTO{
			Code:        p.Code,
			Name:        p.Name,
			Description: p.Description,
			Article:     p.Article,
			Unit:        p.Unit.Name,
			VatCode:     p.VatTax.Code,
			VatTax:      &vatValue,
			GroupCode:   groupCode,
		})
	}

	return response, nil
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

func (svc *productService) limitErrors(skipped []map[string]any, limit int) []map[string]any {
	if len(skipped) > limit {
		return skipped[:limit]
	}
	return skipped
}
