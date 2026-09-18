package repository

import (
	"github.com/victoryus84/gorders/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 1. INTERFAȚA PUBLICĂ (Fără steluță!)
type ContractRepository interface {
	// Contract methods
	CreateContract(contract *models.Contract) error
	UpsertContractBatch(contracts []*models.Contract, batchSize int) error
	FindContractsByClientID(clientID uint) ([]models.Contract, error)
	FindContractByID(id uint) (*models.Contract, error)
	// Contract Address methods
	CreateContractAddress(addr *models.ContractAddress) error
	FindContractAddressByID(id uint) (*models.ContractAddress, error)
}

// 2. STRUCTURA PRIVATĂ (Cu "c" mic)
type contractRepository struct {
	db *gorm.DB
}

// 2. Iată CONSTRUCTORUL de care are nevoie Uber Fx!
func NewContractRepository(db *gorm.DB) ContractRepository {
	return &contractRepository{db: db}
}

// Contract methods
func (rep *contractRepository) CreateContract(contract *models.Contract) error {
	return rep.db.Create(contract).Error
}

func (rep *contractRepository) UpsertContractBatch(contracts []*models.Contract, batchSize int) error {
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

func (rep *contractRepository) FindContractByID(id uint) (*models.Contract, error) {
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

func (rep *contractRepository) FindContractsByClientID(clientID uint) ([]models.Contract, error) {
	var contracts []models.Contract
	err := rep.db.Where("client_id = ?", clientID).Find(&contracts).Error
	return contracts, err
}

func (rep *contractRepository) CreateContractAddress(addr *models.ContractAddress) error {
	return rep.db.Create(addr).Error
}

func (rep *contractRepository) FindContractAddressByID(id uint) (*models.ContractAddress, error) {
	var addr models.ContractAddress
	err := rep.db.First(&addr, id).Error
	// 1. Dacă a apărut o eroare (nu există în DB sau e picat serverul)
	if err != nil {
		return nil, err // Returnăm "mâna goală" (nil) și eroarea
	}
	// 2. Dacă totul e ok, returnăm adresa obiectului plin
	return &addr, nil
}
