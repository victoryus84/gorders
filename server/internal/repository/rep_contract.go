package repository

import (
	"github.com/victoryus84/gorders/internal/models"
)

// Contract methods
func (repository *Repository) CreateContract(contract *models.Contract) error {
	return repository.db.Create(contract).Error
}

func (repository *Repository) FindContractByID(id uint) (*models.Contract, error) {
	var contract models.Contract
	// Încercăm să găsim contractul după ID
	err := repository.db.First(&contract, id).Error

	 // 1. Dacă a apărut o eroare (nu există în DB sau e picat serverul)
    if err != nil {
        return nil, err // Returnăm "mâna goală" (nil) și eroarea
    }
	// 2. Dacă totul e ok, returnăm adresa obiectului plin
	return &contract, nil
}

func (repository *Repository) CreateContractAddress(addr *models.ContractAddress) error {
	return repository.db.Create(addr).Error
}

func (repository *Repository) FindContractAddressByID(id uint) (*models.ContractAddress, error) {
	var addr models.ContractAddress
	err := repository.db.First(&addr, id).Error
	// 1. Dacă a apărut o eroare (nu există în DB sau e picat serverul)
    if err != nil {
        return nil, err // Returnăm "mâna goală" (nil) și eroarea
    }
	// 2. Dacă totul e ok, returnăm adresa obiectului plin
	return &addr, nil
}

func (repository *Repository) FindContractByClientID(clientID uint) ([]models.Contract, error) {
	var contracts []models.Contract
	err := repository.db.Where("client_id = ?", clientID).Find(&contracts).Error
	return contracts, err
}