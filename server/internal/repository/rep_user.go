package repository

import (
	"github.com/victoryus84/gorders/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

// Creează un nou utilizator în baza de date
func (rep *UserRepository) CreateUser(user *models.User) error {
	return rep.db.Create(user).Error
}

// Găsește un utilizator după email
func (rep *UserRepository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := rep.db.Where("email = ?", email).First(&user).Error
	return &user, err
}
