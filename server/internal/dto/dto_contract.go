package dto

type ContractDTO struct {
	Number   string  `json:"number" xml:"number"`
	Name     string  `json:"name" xml:"name" binding:"required"`
	Date     string  `json:"date" xml:"date"` // Format YYYY-MM-DD
	Amount   float64 `json:"amount" xml:"amount"`
	Status   string  `json:"status" xml:"status"`
	FiscalID string  `json:"fiscal_id" xml:"fiscal_id" binding:"required"` // Pentru că 1C trimite Codul Fiscal, nu ID-ul de Client din Postgres
}

type ContractAddressDTO struct {
	Name	   string `json:"name" xml:"name" binding:"required"`
	ContractID uint   `json:"contract_id" xml:"contract_id" binding:"required"`
	Address    string `json:"address" xml:"address" binding:"required"`
	Type       string `json:"type" xml:"type"` // billing, shipping
}
