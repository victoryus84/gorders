package models

import (
	"time"

	"gorm.io/gorm"
)

// ********** Contract - Contract cu clientul **********
type Contract struct {
	gorm.Model
	UUIDModel `gorm:"embedded"`
	Name      string            `gorm:"type:varchar(100);not null"`        // Numele contractului
	Number    *string           `gorm:"type:varchar(50)"`                  // Numărul contractului
	Date      *time.Time        `gorm:"type:date"`                         // Data contractului
	Amount    float64           `gorm:"type:decimal(10,2);not null"`       // Suma contractului
	Status    string            `gorm:"type:varchar(20);not null"`         // Statutul ("active", "closed" etc.)
	ClientID  uint              `gorm:"not null"`                          // Cheie externă către Client
	Client    Client            `gorm:"foreignKey:ClientID;references:ID"` // Clientul
	OwnerID   uint              `gorm:"not null"`                          // ID-ul ownerului (utilizatorului)
	Owner     User              `gorm:"foreignKey:OwnerID;references:ID"`  // Ownerul contractului
	Addresses []ContractAddress `gorm:"foreignKey:ContractID"`             // Adresele asociate contractului
}

// ********** ContractAddress - Adresă asociată contractului **********
type ContractAddress struct {
	gorm.Model
	UUIDModel  `gorm:"embedded"`
	Address    string   `gorm:"type:text;not null"`                  // Adresa
	Type       string   `gorm:"type:varchar(50)"`                    // Tipul adresei ("billing", "shipping" etc.)
	ContractID uint     `gorm:"not null"`                            // Cheie externă către Contract
	Contract   Contract `gorm:"foreignKey:ContractID;references:ID"` // Contractul
	OwnerID    uint     `gorm:"not null"`                            // ID-ul ownerului (utilizatorului)
	Owner      User     `gorm:"foreignKey:OwnerID;references:ID"`    // Ownerul adresei
}
