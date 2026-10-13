package models

import (
	"gorm.io/gorm"
)

// ********** Currency **********
type Currency struct {
	gorm.Model
	UUIDModel 	`gorm:"embedded"`
	Code string `gorm:"type:varchar(10);not null;unique"`
	Name string `gorm:"type:varchar(100);not null"`
	Description string `gorm:"type:text"`
	Symbol string `gorm:"type:varchar(10)"`
	IsBase bool `gorm:"default:false"`
}


