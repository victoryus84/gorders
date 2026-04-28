package repository

import (
	"gorm.io/gorm"
)

// Repository - structura principală pentru acces la baza de date
type Repository struct {
	db *gorm.DB
}

// Creează o nouă instanță de Repository cu conexiunea la DB
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}