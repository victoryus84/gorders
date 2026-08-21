package dto

type ClientDTO struct {
	ID            uint   `json:"id,omitempty" xml:"id,omitempty"`
	ClientTypeID  uint   `json:"client_type" xml:"client_type" binding:"required"`
	Code    	  string `json:"code" xml:"code" binding:"required"`
	Name          string `json:"name" xml:"name" binding:"required"`
	Email         string `json:"email,omitempty" xml:"email,omitempty"`
	Phone         string `json:"phone,omitempty" xml:"phone,omitempty"`
	FiscalAddress string `json:"fiscal_address" xml:"fiscal_address"`
	PostalAddress string `json:"postal_address" xml:"postal_address"`
	FiscalID      string `json:"fiscal_id" xml:"fiscal_id" binding:"required"`
	GroupCode     string `json:"client_group" xml:"client_group" binding:"required"`
}

type ClientGroupDTO struct {
	ID          uint   `json:"id,omitempty" xml:"id,omitempty"`
	Code        string `json:"code" xml:"code" binding:"required"`
	Name        string `json:"name" xml:"name" binding:"required"`
	Description string `json:"description,omitempty" xml:"description,omitempty"`
	ParentCode  string `json:"parent_code,omitempty" xml:"parent_code,omitempty"` 
}

type ClientAddressDTO struct {
	SyncID       string `json:"sync_id"`
	Code		 string `json:"code" xml:"code" binding:"required"` // client code
	Name         string `json:"name"`
	Address      string `json:"address"`
	DeliveryDays int16  `json:"delivery_days"` // Primește direct suma biților din 1C! (ex: 76)
	Type     string `json:"type" xml:"type"` // billing, shipping
}

// Rezultat pentru import masiv
type ImportResult struct {
	Status         string              `json:"status"`
	TotalProcessed int                 `json:"total_processed"`
	TotalSkipped   int                 `json:"total_skipped"`
	ErrorsPreview  []map[string]string `json:"errors_preview"`
	Message        string              `json:"message"`
}