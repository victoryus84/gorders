package repository

import (
	"github.com/victoryus84/gorders/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ContractRepository struct {
	db *gorm.DB
}

// 2. Iată CONSTRUCTORUL de care are nevoie Uber Fx!
func NewContractRepository(db *gorm.DB) *ContractRepository {
	return &ContractRepository{db: db}
}

// Contract methods
func (rep *ContractRepository) CreateContract(contract *models.Contract) error {
	return rep.db.Create(contract).Error
}

func (rep *ContractRepository) UpsertContractBatch(contracts []*models.Contract, batchSize int) error {
	regulaConflict := clause.OnConflict{
		Columns: []clause.Column{
			{Name: "sync_id"}, // <-- Baza de date caută dubluri după asta!
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"name", "number", "client_id", "date", "amount", "status", "updated_at",
		}),
	}

	return rep.db.Clauses(regulaConflict).CreateInBatches(contracts, batchSize).Error
}

func (rep *ContractRepository) FindContractByID(id uint) (*models.Contract, error) {
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

func (rep *ContractRepository) FindContractsByClientID(clientID uint) ([]models.Contract, error) {
	var contracts []models.Contract
	err := rep.db.Where("client_id = ?", clientID).Find(&contracts).Error
	return contracts, err
}

func (rep *ContractRepository) CreateContractAddress(addr *models.ContractAddress) error {
	return rep.db.Create(addr).Error
}

func (rep *ContractRepository) FindContractAddressByID(id uint) (*models.ContractAddress, error) {
	var addr models.ContractAddress
	err := rep.db.First(&addr, id).Error
	// 1. Dacă a apărut o eroare (nu există în DB sau e picat serverul)
	if err != nil {
		return nil, err // Returnăm "mâna goală" (nil) și eroarea
	}
	// 2. Dacă totul e ok, returnăm adresa obiectului plin
	return &addr, nil
}
