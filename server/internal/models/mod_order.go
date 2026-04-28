package models

import (
	"gorm.io/gorm"
)

// ********** Order - Comandă **********
type Order struct {
	gorm.Model
	UUIDModel   `gorm:"embedded"`
	OwnerID     uint        `gorm:"not null"`                            // ID-ul ownerului (utilizatorului)
	Owner       User        `gorm:"foreignKey:OwnerID;references:ID"`    // Ownerul comenzii
	ClientID    uint        `gorm:"not null"`                            // ID-ul clientului (cheie externă)
	Client      Client      `gorm:"foreignKey:ClientID;references:ID"`   // Clientul care a plasat comanda
	PriceTypeID uint        `gorm:"not null"`                            // ID-ul tipului de preț (cheie externă)
	PriceType   PriceType   `gorm:"foreignKey:PriceTypeID"`              // Tipul de preț al comenzii
	ContractID  uint        `gorm:"not null"`                            // ID-ul contractului (cheie externă)
	Contract    Contract    `gorm:"foreignKey:ContractID;references:ID"` // Contractul asociat comenzii
	TotalPrice  float64     `gorm:"type:decimal(10,2);not null"`         // Suma totală a comenzii
	Status      string      `gorm:"type:varchar(20);not null"`           // Statusul comenzii
	OrderItems  []OrderItem `gorm:"foreignKey:OrderID"`                  // Pozițiile comenzii
}

// ****************************************************

// ********** OrderItem - Poziție comandă **********
type OrderItem struct {
	gorm.Model
	UUIDModel   `gorm:"embedded"`
	OrderID     uint    `gorm:"not null"`                           // ID-ul comenzii
	ProductID   uint    `gorm:"not null"`                           // ID-ul produsului
	Product     Product `gorm:"foreignKey:ProductID;references:ID"` // Produsul asociat poziției
	Quantity    float64 `gorm:"type:decimal(10,3);not null"`        // Cantitatea
	Price       float64 `gorm:"type:decimal(10,2);not null"`        // Prețul unitar la momentul comenzii
	UnitID      uint    `gorm:"not null"`                           // ID-ul unității de măsură
	Unit        Unit    `gorm:"foreignKey:UnitID;references:ID"`    // Unitatea de măsură asociată poziției
	UnitName    string  `gorm:"type:varchar(20)"`                   // Stocăm "KG" sau "BUC"
	VatTaxID    uint    `gorm:"not null"`                           // ID-ul taxei VAT
	VatTax      VatTax  `gorm:"foreignKey:VatTaxID;references:ID"`  // Taxa VAT asociată poziției
	VatRate     float64 `gorm:"type:decimal(10,2);not null"`        // Rata TVA-ului (preluată din VatTax)
	Summ        float64 `gorm:"type:decimal(10,2);not null"`        // Suma totală pentru poziție (Price * Quantity)
	VatSumm     float64 `gorm:"type:decimal(10,2);not null"`        // Valoarea TVA-ului în bani
	SummWithVat float64 `gorm:"type:decimal(10,2);not null"`        // Suma totală pentru poziție cu TVA (Summ + VAT)
}
