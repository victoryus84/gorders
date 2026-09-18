package repository

import (
	"github.com/victoryus84/gorders/internal/models"
	"gorm.io/gorm"
)

// 1. INTERFAȚA PUBLICĂ (Fără steluță!)
type OrderRepository interface {
	// Order methods
	CreateOrder(order *models.Order) error
	FindOrdersByUserID(userID uint) ([]models.Order, error)
	FindOrderByID(id uint) (*models.Order, error)
}

// 2. STRUCTURA PRIVATĂ (Cu "c" mic)
type orderRepository struct {
	db *gorm.DB
}

// 2. Iată CONSTRUCTORUL de care are nevoie Uber Fx!
func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}
// Order methods
func (rep *orderRepository) CreateOrder(order *models.Order) error {
	return rep.db.Create(order).Error
}
func (rep *orderRepository) FindOrdersByUserID(userID uint) ([]models.Order, error) {
	var orders []models.Order
	err := rep.db.Preload("OrderItems").Where("user_id = ?", userID).Find(&orders).Error
	return orders, err
}
func (rep *orderRepository) FindOrderByID(id uint) (*models.Order, error) {
	var order models.Order
	err := rep.db.Preload("OrderItems").First(&order, id).Error
	return &order, err
}
