package dto

type ProductDTO struct {
	SyncID string  `json:"sync_id" xml:"sync_id"`
	Number string  `json:"number" xml:"number"`
	Name   string  `json:"name" xml:"name" binding:"required"`
	Date   string  `json:"date" xml:"date"` // Format YYYY-MM-DD
	Amount float64 `json:"amount" xml:"amount"`
	Status string  `json:"status" xml:"status"`
	Code   string  `json:"code" xml:"code" binding:"required"` // client code
}
