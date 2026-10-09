package dto

// DTO pentru Valute
type ValuteDTO struct {
	SyncID string `json:"sync_id" xml:"sync_id"` // ID opțional pentru sincronizare
	Name string `json:"name" xml:"name" binding:"required"`
	Description string `json:"description" xml:"description"`
	Code string `json:"code" xml:"code" binding:"required"`
	Symbol string `json:"symbol" xml:"symbol"`
	IsBase bool `json:"isbase" xml:"isbase" binding:"required"`
}

