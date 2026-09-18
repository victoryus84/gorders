package service

import (
	"github.com/victoryus84/gorders/internal/dto"
	"github.com/victoryus84/gorders/internal/logger"
	"github.com/victoryus84/gorders/internal/models"
	"github.com/victoryus84/gorders/internal/repository"
)

// 1. INTERFAȚA PUBLICĂ - Asta e tot ce văd Handlerele sau alte module
type ProductService interface {
	ProcessProductImport(dtos []dto.ProductDTO) map[string]interface{}
	ProcessProductGroupImport(dtos []dto.ProductGroupDTO) map[string]interface{}
	GetFirst1000Products() ([]models.Product, error)
}

// 2. STRUCTURA PRIVATĂ - Aici ținem "uneltele" (Repozitoarele dedicate)
type productService struct {
	productRepository repository.ProductRepository
	unitRepository    repository.UnitRepository
	vatRepository     repository.VatTaxRepository
}

// 3. CONSTRUCTORUL PENTRU UBER FX
// Primește repozitoarele tale specifice și returnează Interfața
func NewProductService(
	pr repository.ProductRepository,
	ur repository.UnitRepository,
	vr repository.VatTaxRepository,
) ProductService {
	return &productService{
		productRepository: pr,
		unitRepository:    ur,
		vatRepository:     vr,
	}
}

// 4. IMPLEMENTAREA LOGICII DE SINCRONIZARE
func (svc *productService) ProcessProductImport(dtos []dto.ProductDTO) map[string]interface{} {
	if len(dtos) == 0 {
		return map[string]interface{}{"received": 0, "inserted": 0, "errors": 0}
	}

	// Încărcăm dicționarele folosind metodele tale clare pe care le vom scrie în repozitoarele lor
	unitMap, err := svc.unitRepository.GetUnitIDMap() 
	if err != nil {
		logger.LogError("❌ Lipsă dicționar Unități", err)
		return map[string]interface{}{"error": "Lipsește dicționarul de Unități"}
	}

	vatMap, err := svc.vatRepository.GetVatIDMap()
	if err != nil {
		logger.LogError("❌ Lipsă dicționar TVA", err)
		return map[string]interface{}{"error": "Lipsește dicționarul de TVA"}
	}

	// Traducem și pregătim modelele
	items := make([]models.Product, 0, len(dtos))
	for _, input := range dtos {
		item := models.Product{
			Code:        input.Code,
			Name:        input.Name,
			Description: input.Description,
			Article:     input.Article,
			UnitID:      unitMap[input.Unit],    // Tradus instant
			VatTaxID:    vatMap[input.VatCode], // Tradus instant
		}
		items = append(items, item)
	}

	// Salvăm lotul folosind metoda ta dedicată
	err = svc.productRepository.UpsertProductsBatch(items, 500)
	if err != nil {
		logger.LogError("❌ Eroare la salvarea produselor", err)
		return map[string]interface{}{"received": len(dtos), "inserted": 0, "errors": len(dtos)}
	}

	return map[string]interface{}{"received": len(dtos), "inserted": len(items), "errors": 0}
}

// 1. Funcția care aduce primii 1000 de produse
func (svc *productService) GetFirst1000Products() ([]models.Product, error) {
	// Trimitem cererea mai departe către repozitoriu
	return svc.productRepository.GetFirst1000Products()
}

// 2. Funcția pentru Grupe de Produse (Schelet ca să nu urle compilatorul)
func (svc *productService) ProcessProductGroupImport(dtos []dto.ProductGroupDTO) map[string]interface{} {
	// TODO: Aici vom face maparea din DTO în Model și vom chema un UpsertProductGroupsBatch
	// Momentan returnăm un răspuns gol ca să compileze perfect aplicația.
	return map[string]interface{}{
		"status": "funcția pentru grupe este în construcție",
	}
}