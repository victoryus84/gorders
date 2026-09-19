package models

import (
	"gorm.io/gorm"
)

// ********** Product - Produs **********
type Product struct {
	gorm.Model
	UUIDModel      `gorm:"embedded"`
	Code           string       `gorm:"type:varchar(15);unique;not null"`        // Codul produsului
	Name           string       `gorm:"type:varchar(100);not null"`              // Numele produsului
	Description    string       `gorm:"type:text"`                               // Descrierea produsului
	Price          float64      `gorm:"type:decimal(10,2);default:0.0"`          // Prețul produsului
	ProductGroup   ProductGroup `gorm:"foreignKey:ProductGroupID;references:ID"` // Grupa de produse din care face parte
	ProductGroupID *uint        `gorm:"column:product_group_id"`                 // ID-ul grupei de produse
	Unit           Unit         `gorm:"foreignKey:UnitID;references:ID"`         // Unitatea de măsură a produsului
	UnitID         uint         `gorm:"not null"`                                // ID-ul unității de măsură
	VatTax         VatTax       `gorm:"foreignKey:VatTaxID;references:ID"`       // Taxa VAT a produsului
	VatTaxID       uint         `gorm:"not null"`                                // ID-ul taxei VAT
	Article        string       `gorm:"type:varchar(100)"`                       // Articolul produsului (opțional)
}

// ********** ProductGroup - Grupa de Produse **********
type ProductGroup struct {
	gorm.Model
	UUIDModel   `gorm:"embedded"`
	Code        string    `gorm:"type:varchar(15);unique;not null"`  // Codul grupei (ex: "001", "ELEC")
	Name        string    `gorm:"type:varchar(100);not null"` // Numele grupei (ex: "Băuturi", "Electronice")
	Description string    `gorm:"type:text"`                         // Descrierea grupei
	ParentID    *uint 
	Parent      *ProductGroup `gorm:"foreignKey:ParentID"`            // Legătură către grupa părinte (dacă există)
	Children    []ProductGroup `gorm:"foreignKey:ParentID"`            // Legătură către grupele copil (dacă există)
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
