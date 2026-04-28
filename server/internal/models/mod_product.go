package models

import (
	"gorm.io/gorm"
)

// ********** Product - Produs **********
type Product struct {
	gorm.Model
	UUIDModel      `gorm:"embedded"`
	Name           string       `gorm:"type:varchar(100);not null"`              // Numele produsului
	Price          float64      `gorm:"type:decimal(10,2);default:0.0"`          // Prețul produsului
	Description    string       `gorm:"type:text"`                               // Descrierea produsului
	ProductGroupID uint         `gorm:"not null"`                                // ID-ul grupei de produse
	ProductGroup   ProductGroup `gorm:"foreignKey:ProductGroupID;references:ID"` // Grupa de produse din care face parte
	UnitID         uint         `gorm:"not null"`                                // ID-ul unității de măsură
	Unit           Unit         `gorm:"foreignKey:UnitID;references:ID"`         // Unitatea de măsură a produsului
	VatTaxID       uint         `gorm:"not null"`                                // ID-ul taxei VAT
	VatTax         VatTax       `gorm:"foreignKey:VatTaxID;references:ID"`       // Taxa VAT a produsului
}

// ********** ProductGroup - Grupa de Produse **********
type ProductGroup struct {
	gorm.Model
	UUIDModel   `gorm:"embedded"`
	Name        string    `gorm:"type:varchar(100);not null;unique"` // Numele grupei (ex: "Băuturi", "Electronice")
	Description string    `gorm:"type:text"`                         // Descrierea grupei
	Products    []Product `gorm:"foreignKey:ProductGroupID"`         // O grupă are mai multe produse
}

// ********** Price type of products - Tipuri de pret **********
type PriceType struct {
	gorm.Model
	UUIDModel   `gorm:"embedded"`
	Name        string `gorm:"type:varchar(50);not null"` // Numele tipului de preț (ex: "Cu amamuntul", "En-gros")
	Description string `gorm:"type:text"`                 // Descrierea tipiului de preț
}

// ********** Price of products - Preturi producte **********
type PriceProduct struct {
	gorm.Model
	UUIDModel   `gorm:"embedded"`
	ProductID   uint      `gorm:"not null"`                             // Cheie externă către Product
	Product     Product   `gorm:"foreignKey:ProductID;references:ID"`   // Produsul
	PriceTypeID uint      `gorm:"not null"`                             // Cheie externă către PriceType
	PriceType   PriceType `gorm:"foreignKey:PriceTypeID;references:ID"` // Tipul de preț
	Price       float64   `gorm:"type:decimal(10,2);not null"`          // Prețul pentru acest tip
}
