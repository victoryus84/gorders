package dto

type ClientDTO struct {
	ID            uint   `json:"id,omitempty" xml:"id,omitempty"`
	ClientTypeID  uint   `json:"client_type" xml:"client_type" binding:"required"`
	Name          string `json:"name" xml:"name" binding:"required"`
	FiscalID      string `json:"fiscal_id" xml:"fiscal_id" binding:"required"`
	Email         string `json:"email,omitempty" xml:"email,omitempty"`
	Phone         string `json:"phone,omitempty" xml:"phone,omitempty"`
	FiscalAddress string `json:"fiscal_address" xml:"fiscal_address"`
	PostalAddress string `json:"postal_address" xml:"postal_address"`
}

type ClientAddressDTO struct {
	ID       uint   `json:"id,omitempty" xml:"id,omitempty"`
	FiscalID string `json:"fiscal_id" xml:"fiscal_id" binding:"required"`
	Name     string `json:"name" xml:"name" binding:"required"`
	Address  string `json:"address" xml:"address" binding:"required"`
	Type     string `json:"type" xml:"type"` // billing, shipping
}

// Rezultat pentru import masiv
type ImportResult struct {
	Status        string              `json:"status"`
	TotalCreated  int                 `json:"total_created"`
	TotalSkipped  int                 `json:"total_skipped"`
	ErrorsPreview []map[string]string `json:"errors_preview"`
	Message       string              `json:"message"`
}