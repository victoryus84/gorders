package models

import (
	"gorm.io/gorm"
)

// ********** Client - Client (beneficiar) **********
type Client struct {
	gorm.Model
	UUIDModel     `gorm:"embedded"`
	ClientTypeID  uint            `gorm:"not null"`                         // Foreign key to ClientType
	ClientType    ClientType      `gorm:"foreignKey:ClientTypeID;not null"` // Tipul clientului ("individual", "company", etc.)
	Name          string          `gorm:"type:varchar(100);not null"`       // Numele clientului
	FiscalID      string          `gorm:"type:varchar(15);unique;not null"` // Codul fiscal al clientului (unic)
	Email         *string         `gorm:"type:varchar(100);unique"`         // Email-ul clientului (unic)
	Phone         string          `gorm:"type:varchar(50)"`                 // Telefonul clientului
	FiscalAddress string          `gorm:"type:text"`                        // Adresa fiscală a clientului
	PostalAddress string          `gorm:"type:text"`                        // Adresa postala a clientului
	Contracts     []Contract      `gorm:"foreignKey:ClientID"`              // Contractele clientului
	Addresses     []ClientAddress `gorm:"foreignKey:ClientID"`              // Adresele asociate clientului
}

// ********** ClientAddress - Adresă asociată clientului **********
type ClientAddress struct {
	gorm.Model
	UUIDModel `gorm:"embedded"`
	Name      string  `gorm:"type:varchar(100);not null"`        // Numele adresei
	Address   *string `gorm:"type:text"`                         // Adresa
	Type      string  `gorm:"type:varchar(50)"`                  // Tipul adresei ("billing", "shipping" etc.)
	ClientID  uint    `gorm:"not null"`                          // Cheie externă către Client
	Client    Client  `gorm:"foreignKey:ClientID;references:ID"` // Clientul
	OwnerID   uint    `gorm:"not null"`                          // ID-ul ownerului (utilizatorului)
	Owner     User    `gorm:"foreignKey:OwnerID;references:ID"`  // Ownerul adresei
}

// ********** Client - Client (beneficiar) **********
type ClientType struct {
	gorm.Model
	UUIDModel `gorm:"embedded"`
	Name      string `gorm:"type:varchar(20);not null"` // Tipul clientului ("individual", "company", etc.)
}
