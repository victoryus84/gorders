package models

import (
	"gorm.io/gorm"
)

// ********** Channel - Canal de vânzări **********
type Channel struct {
	gorm.Model
	UUIDModel   `gorm:"embedded"`
	Name        string `gorm:"type:varchar(100);not null"` // Numele canalului
	Description string `gorm:"type:text"`                  // Descrierea canalului
	Users       []User `gorm:"many2many:user_channels;"`   // Utilizatorii care aparțin acestui canal
}
