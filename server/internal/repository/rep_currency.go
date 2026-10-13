package repository

import (
	"github.com/victoryus84/gorders/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 1. INTERFAȚA
type CurrencyRepository interface {
	UpsertBatch(items []models.Currency, batchSize int) error
	FindAll() ([]models.Currency, error)
}

// 2. STRUCTURA PRIVATĂ
type currencyRepository struct {
	db *gorm.DB
}

// 3. CONSTRUCTORUL UBER FX
func NewCurrencyRepository(db *gorm.DB) CurrencyRepository {
	return &currencyRepository{db: db}
}

// 4. METODELE PRINCIPALE
func (r *currencyRepository) UpsertBatch(items []models.Currency, batchSize int) error {
	return r.db.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(items, batchSize).Error
}

func (r *currencyRepository) FindAll() ([]models.Currency, error) {
	var items []models.Currency
	err := r.db.Find(&items).Error
	return items, err
}


