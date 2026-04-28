package service

import (
	"github.com/victoryus84/gorders/internal/dto"
	"github.com/victoryus84/gorders/internal/models"
)

type ContractRepository interface {
	CreateContract(contract *models.Contract) error
	FindContractByID(id uint) (*models.Contract, error)
	FindContractsByClientID(clientID uint) ([]models.Contract, error)
	CreateContractAddress(address *models.ContractAddress) error
	FindContractAddressByID(id uint) (*models.ContractAddress, error)
}

type ContractService interface {
	CreateContract(requests []dto.ContractDTO, ownerID uint) dto.ImportResult
	FindContractByID(id uint) (*dto.ContractDTO, error)
	FindContractByClientID(clientID uint) ([]dto.ContractDTO, error)
	CreateContractAddress(requests []dto.ContractAddressDTO, ownerID uint) dto.ImportResult
	FindContractAddressByID(id uint) (*dto.ContractAddressDTO, error)
}

type contractService struct {
	repo ContractRepository
}		

func NewContractService(repo ContractRepository) ContractService {
	return &contractService{repo: repo}
}	

func (s *contractService) ProcessContractImport(requests []dto.ContractDTO) dto.ImportResult {
	created := make([]*models.Contract, 0)
	skipped := make([]map[string]string, 0)

	for _, req := range requests {
		// Process each contract request
	}

}