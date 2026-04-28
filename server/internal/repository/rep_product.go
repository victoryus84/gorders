package repository

import (
	"github.com/victoryus84/gorders/internal/models"
)

// Product methods
func (repository *Repository) CreateProduct(product *models.Product) error {
	return repository.db.Create(product).Error
}
func (repository *Repository) FindProductGroupByID(id uint) (*models.ProductGroup, error) {
	var group models.ProductGroup
	err := repository.db.First(&group, id).Error
	return &group, err
}
func (repository *Repository) FindProductByID(id uint) (*models.Product, error) {
	var product models.Product
	err := repository.db.First(&product, id).Error
	return &product, err
}
func (repository *Repository) FindVatTaxByID(id uint) (*models.VatTax, error) {
	var vatTax models.VatTax
	err := repository.db.First(&vatTax, id).Error
	return &vatTax, err
}
func (repository *Repository) FindUnitByID(id uint) (*models.Unit, error) {
	var unit models.Unit
	err := repository.db.First(&unit, id).Error
	return &unit, err
}