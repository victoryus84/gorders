package repository

import (
	"github.com/victoryus84/gorders/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 1. INTERFAȚA
type ValuteRepository interface {
	UpsertBatch(items []models.Valute, batchSize int) error
	FindAll() ([]models.Valute, error)
}

// 2. STRUCTURA PRIVATĂ
type valuteRepository struct {
	db *gorm.DB
}

// 3. CONSTRUCTORUL UBER FX
func NewValuteRepository(db *gorm.DB) ValuteRepository {
	return &valuteRepository{db: db}
}

// 4. METODELE PRINCIPALE
func (r *valuteRepository) UpsertBatch(items []models.Valute, batchSize int) error {
	return r.db.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(items, batchSize).Error
}

func (r *valuteRepository) FindAll() ([]models.Valute, error) {
	var items []models.Valute
	err := r.db.Find(&items).Error
	return items, err
}


