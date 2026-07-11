package repository

import (
	"github.com/victoryus84/gorders/internal/models"
	"gorm.io/gorm/clause"
)

// Contract methods
func (rep *Repository) CreateContract(contract *models.Contract) error {
	return rep.db.Create(contract).Error
}

func (rep *Repository) UpsertContractBatch(contracts []*models.Contract, batchSize int) error {
    regulaConflict := clause.OnConflict{
        Columns:   []clause.Column{{Name: "client_id"}, {Name: "number"}}, // Cheia unică pentru conflict
        DoUpdates: clause.AssignmentColumns([]string{
            "name", "number", "date", "amount", "status", "client_id", "owner_id", "updated_at",
        }),
    }

    // Aici e magia: CreateInBatches
    return rep.db.Clauses(regulaConflict).CreateInBatches(contracts, batchSize).Error
}

func (rep *Repository) FindContractByID(id uint) (*models.Contract, error) {
	var contract models.Contract
	// Încercăm să găsim contractul după ID
	err := rep.db.First(&contract, id).Error

	 // 1. Dacă a apărut o eroare (nu există în DB sau e picat serverul)
    if err != nil {
        return nil, err // Returnăm "mâna goală" (nil) și eroarea
    }
	// 2. Dacă totul e ok, returnăm adresa obiectului plin
	return &contract, nil
}
func (rep *Repository) FindContractsByClientID(clientID uint) ([]models.Contract, error) {
	var contracts []models.Contract
	err := rep.db.Where("client_id = ?", clientID).Find(&contracts).Error
	return contracts, err
}
func (rep *Repository) CreateContractAddress(addr *models.ContractAddress) error {
	return rep.db.Create(addr).Error
}
func (rep *Repository) FindContractAddressByID(id uint) (*models.ContractAddress, error) {
	var addr models.ContractAddress
	err := rep.db.First(&addr, id).Error
	// 1. Dacă a apărut o eroare (nu există în DB sau e picat serverul)
    if err != nil {
        return nil, err // Returnăm "mâna goală" (nil) și eroarea
    }
	// 2. Dacă totul e ok, returnăm adresa obiectului plin
	return &addr, nil
}
