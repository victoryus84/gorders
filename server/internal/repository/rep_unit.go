package repository

import (
	"github.com/victoryus84/gorders/internal/models"
	"gorm.io/gorm"
)

// 1. Interfața pe care o va cere ProductService
type UnitRepository interface {
	GetUnitIDMap() (map[string]uint, error)
}

// 2. Structura privată care ține conexiunea
type unitRepository struct {
	db *gorm.DB
}

// 3. Constructorul pentru Uber Fx
func NewUnitRepository(db *gorm.DB) UnitRepository {
	return &unitRepository{db: db}
}

// 4. Logica care aduce dicționarul ultra-rapid
func (r *unitRepository) GetUnitIDMap() (map[string]uint, error) {
	// Facem un query rapid doar pentru ID și Name
	rows, err := r.db.Model(&models.Unit{}).Select("id, name").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Creăm dicționarul în memorie
	resultMap := make(map[string]uint)
	
	for rows.Next() {
		var id uint
		var name string
		
		// Citim valorile din baza de date
		if err := rows.Scan(&id, &name); err != nil {
			continue 
		}
		resultMap[name] = id
	}
	
	return resultMap, nil
}