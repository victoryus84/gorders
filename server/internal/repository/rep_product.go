package repository

import (
	"github.com/victoryus84/gorders/internal/models"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

// Product methods
func (rep *ProductRepository) CreateProduct(product *models.Product) error {
	return rep.db.Create(product).Error
}
func (rep *ProductRepository) FindProductGroupByID(id uint) (*models.ProductGroup, error) {
	var group models.ProductGroup
	err := rep.db.First(&group, id).Error
	return &group, err
}
func (rep *ProductRepository) FindProductByID(id uint) (*models.Product, error) {
	var product models.Product
	err := rep.db.First(&product, id).Error
	return &product, err
}
func (rep *ProductRepository) FindVatTaxByID(id uint) (*models.VatTax, error) {
	var vatTax models.VatTax
	err := rep.db.First(&vatTax, id).Error
	return &vatTax, err
}
func (rep *ProductRepository) FindUnitByID(id uint) (*models.Unit, error) {
	var unit models.Unit
	err := rep.db.First(&unit, id).Error
	return &unit, err
}
