package repository

import (
	"github.com/victoryus84/gorders/internal/models"
	"gorm.io/gorm"
)

// 1. Interfața publică
type UserRepository interface {
    CreateUser(user *models.User) error
    FindUserByEmail(email string) (*models.User, error)
}

// 2. Structura privată (cu literă mică)
type userRepository struct {
    db *gorm.DB // sau ce folosești tu pentru DB
}

// 3. Constructorul trebuie să returneze interfața!
func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}

// Creează un nou utilizator în baza de date
func (rep *userRepository) CreateUser(user *models.User) error {
	return rep.db.Create(user).Error
}

// Găsește un utilizator după email
func (rep *userRepository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := rep.db.Where("email = ?", email).First(&user).Error
	return &user, err
}
