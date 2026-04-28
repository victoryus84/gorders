package repository

import (
	"github.com/victoryus84/gorders/internal/models"
)

// Creează un nou utilizator în baza de date
func (repository *Repository) CreateUser(user *models.User) error {
	return repository.db.Create(user).Error
}

// Găsește un utilizator după email
func (repository *Repository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := repository.db.Where("email = ?", email).First(&user).Error
	return &user, err
}
