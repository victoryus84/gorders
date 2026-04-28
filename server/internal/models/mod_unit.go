package models

import (
	"gorm.io/gorm"
)

// ********** Units of Measurement - Unitati de măsură **********
type Unit struct {
	gorm.Model
	UUIDModel   `gorm:"embedded"`
	Name        string  `gorm:"type:varchar(50);not null"`               // Numele unității de măsură (ex: "buc", "kg")
	Description string  `gorm:"type:text"`                               // Descrierea unității de măsură
	Coefficient float64 `gorm:"type:decimal(10,4);not null;default:1.0"` // Coeficient de conversie față de unitatea de bază

}
