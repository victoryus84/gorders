package models

import (
	"gorm.io/gorm"
)

// ********** VatRate - TVA **********
type VatTax struct {
	gorm.Model
	UUIDModel   `gorm:"embedded"`
	Name        string  `gorm:"type:varchar(100);not null"`  // Numele taxei
	Rate        float64 `gorm:"type:decimal(10,2);not null"` // Rata taxei
	Description string  `gorm:"type:text"`                   // Descrierea taxei
}

// ********** IncomeTax - Taxă pe venit **********
type IncomeTax struct {
	gorm.Model
	UUIDModel   `gorm:"embedded"`
	Name        string  `gorm:"type:varchar(100);not null"`  // Numele taxei
	Rate        float64 `gorm:"type:decimal(10,2);not null"` // Rata taxei
	Description string  `gorm:"type:text"`                   // Descrierea taxei
}
