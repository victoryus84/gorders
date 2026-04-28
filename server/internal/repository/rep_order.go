package repository

import (
	"github.com/victoryus84/gorders/internal/models"
)

// Order methods
func (repository *Repository) CreateOrder(order *models.Order) error {
	return repository.db.Create(order).Error
}
func (repository *Repository) FindOrdersByUserID(userID uint) ([]models.Order, error) {
	var orders []models.Order
	err := repository.db.Preload("OrderItems").Where("user_id = ?", userID).Find(&orders).Error
	return orders, err
}
func (repository *Repository) FindOrderByID(id uint) (*models.Order, error) {
	var order models.Order
	err := repository.db.Preload("OrderItems").First(&order, id).Error
	return &order, err
}
