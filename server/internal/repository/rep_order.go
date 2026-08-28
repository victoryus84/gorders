package repository

import (
	"github.com/victoryus84/gorders/internal/models"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

// Order methods
func (rep *OrderRepository) CreateOrder(order *models.Order) error {
	return rep.db.Create(order).Error
}
func (rep *OrderRepository) FindOrdersByUserID(userID uint) ([]models.Order, error) {
	var orders []models.Order
	err := rep.db.Preload("OrderItems").Where("user_id = ?", userID).Find(&orders).Error
	return orders, err
}
func (rep *OrderRepository) FindOrderByID(id uint) (*models.Order, error) {
	var order models.Order
	err := rep.db.Preload("OrderItems").First(&order, id).Error
	return &order, err
}
