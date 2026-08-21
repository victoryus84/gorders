package models

import (
	"gorm.io/gorm"
)

// ********** Client - Client (beneficiar) **********
type Client struct {
	gorm.Model
	UUIDModel     `gorm:"embedded"`
	ClientGroupID *uint           `gorm:"column:client_group_id"`           // Foreign key to ClientGroup
	ClientTypeID  uint            `gorm:"not null"`                         // Foreign key to ClientType
	Code    	  string          `gorm:"type:varchar(15);unique;not null"` // Codul clientului (unic)
	Name          string          `gorm:"type:varchar(100);not null"`       // Numele clientului
	Description   string          `gorm:"type:text"`                        // Descrierea clientului
	FiscalID      string          `gorm:"type:varchar(15);unique;not null"` // Codul fiscal al clientului (unic)
	Email         *string         `gorm:"type:varchar(100);unique"`         // Email-ul clientului (unic)
	Phone         string          `gorm:"type:varchar(50)"`                 // Telefonul clientului
	FiscalAddress string          `gorm:"type:text"`                        // Adresa fiscală a clientului
	PostalAddress string          `gorm:"type:text"`                        // Adresa postala a clientului
	ClientGroup   ClientGroup     `gorm:"foreignKey:ClientGroupID;not null"` // Grupa din care face parte clientul
	ClientType    ClientType      `gorm:"foreignKey:ClientTypeID;not null"` // Tipul clientului ("individual", "company", etc.)
	Contracts     []Contract      `gorm:"foreignKey:ClientID"`              // Contractele clientului
	Addresses     []ClientAddress `gorm:"foreignKey:ClientID"`              // Adresele asociate clientului
}

// ********** ClientGroup - Grupa de Clienti **********
type ClientGroup struct {
	gorm.Model
	UUIDModel   `gorm:"embedded"`
	Code    	string    `gorm:"type:varchar(15);unique;not null"`  // Codul grupei (unic)
	Name        string    `gorm:"type:varchar(100);not null;unique"` // Numele grupei (ex: "Băuturi", "Electronice")
	Description string    `gorm:"type:text"`                         // Descrierea grupei
	ParentID    *uint 
	Parent      *ClientGroup `gorm:"foreignKey:ParentID"`            // Legătură către grupa părinte (dacă există)
	Children    []ClientGroup `gorm:"foreignKey:ParentID"`            // Legătură către grupele copil (dacă există)
	Clients    []Client   `gorm:"foreignKey:ClientGroupID"`          // O grupă are mai mulți clienți
}
// ********** ClientAddress - Adresă asociată clientului **********
type ClientAddress struct {
	gorm.Model
	UUIDModel `gorm:"embedded"`
	SyncID      string  `gorm:"type:varchar(100);uniqueIndex;not null"` // SyncID adresei (unic, pentru sincronizare cu 1C)
	Name        string  `gorm:"type:varchar(100);not null"`             // Numele adresei
	Address     *string `gorm:"type:text"`                              // Adresa
	Type        string  `gorm:"type:varchar(50)"`                       // Tipul adresei ("billing", "shipping" etc.)
	Description *string  `gorm:"type:text"`                             // Descrierea adresei
	DeliveryDays int16  `gorm:"default:0;not null"` 					// Aici ținem bitmask-ul pentru zilele de livrare (0-6 pentru L-D)
	ClientID    uint    `gorm:"not null"`                               // Cheie externă către Client
	Client      Client  `gorm:"foreignKey:ClientID;references:ID"`      // Clientul
	OwnerID     uint    `gorm:"not null"`                               // ID-ul ownerului (utilizatorului)
	Owner       User    `gorm:"foreignKey:OwnerID;references:ID"`       // Ownerul adresei
}

// ********** Client - Client (beneficiar) **********
type ClientType struct {
	gorm.Model
	UUIDModel `gorm:"embedded"`
	Name      string `gorm:"type:varchar(20);not null"` // Tipul clientului ("individual", "company", etc.)
}
