package repository

import (
	"github.com/victoryus84/gorders/internal/models"
	"gorm.io/gorm"
)

// 1. Interfața
type VatTaxRepository interface {
	GetVatIDMap() (map[string]uint, error)
}

// 2. Structura privată
type vatTaxRepository struct {
	db *gorm.DB
}

// 3. Constructorul pentru Uber Fx
func NewVatTaxRepository(db *gorm.DB) VatTaxRepository {
	return &vatTaxRepository{db: db}
}

// 4. Logica pentru dicționar
func (r *vatTaxRepository) GetVatIDMap() (map[string]uint, error) {
	// Selectăm `id` și `code` (identificatorul unic pentru sincronizare)
	rows, err := r.db.Model(&models.VatTax{}).Select("id, code").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resultMap := make(map[string]uint)
	
	for rows.Next() {
		var id uint
		var code string
		
		// Scanăm id-ul și codul
		if err := rows.Scan(&id, &code); err != nil {
			continue // Ignorăm rândurile cu erori la citire
		}
		
		// Un strat extra de siguranță: nu adăugăm în map codurile goale
		// (în caz că ai date vechi care încă nu au primit codul la migrare)
		if code != "" {
			resultMap[code] = id
		}
	}
	
	return resultMap, nil
}