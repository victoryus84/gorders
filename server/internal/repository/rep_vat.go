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
	// Aici citim id și description (codul tău din 1C)
	rows, err := r.db.Model(&models.VatTax{}).Select("id, description").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resultMap := make(map[string]uint)
	
	for rows.Next() {
		var id uint
		var description string
		
		if err := rows.Scan(&id, &description); err != nil {
			continue 
		}
		resultMap[description] = id
	}
	
	return resultMap, nil
}