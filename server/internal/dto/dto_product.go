package dto

type ProductDTO struct {
	SyncID 		string `json:"sync_id" xml:"sync_id"` // optional field for synchronization purposes
	Code   		string `json:"code" xml:"code" binding:"required"` // product code
	Name   		string `json:"name" xml:"name" binding:"required"`
	Description string `json:"description,omitempty" xml:"description,omitempty"`
	Article 	string `json:"article,omitempty" xml:"article,omitempty"`
	Unit        string `json:"unit" xml:"unit" binding:"required"`
	VatCode 	string `json:"vat_code" xml:"vat_code" binding:"required"`
	VatTax	    *float64 `json:"vat_tax" xml:"vat_tax" binding:"required"`
	GroupCode   string `json:"group" xml:"group" binding:"required"`
}

type ProductGroupDTO struct {
	SyncID      string  `json:"sync_id" xml:"sync_id"`
	Code        string `json:"code" xml:"code" binding:"required"`
	Name        string `json:"name" xml:"name" binding:"required"`
	Description string `json:"description,omitempty" xml:"description,omitempty"`
	ParentCode  string `json:"parent_code,omitempty" xml:"parent_code,omitempty"`
}