package service

import (
	"github.com/victoryus84/gorders/internal/models"
	"github.com/victoryus84/gorders/internal/repository"
)

// 1. INTERFAȚA PUBLICĂ
type CurrencyService interface {
	FindAll() ([]models.Currency, error)
}

// 2. STRUCTURA PRIVATĂ
type currencyService struct {
	rep repository.CurrencyRepository
}

func NewCurrencyService(rep repository.CurrencyRepository) CurrencyService {
	return &currencyService{rep: rep}
}

func (svc *currencyService) FindAll() ([]models.Currency, error) {
	return svc.rep.FindAll()
}
