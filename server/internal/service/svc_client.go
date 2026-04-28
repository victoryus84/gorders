package service

import (
	"strings"

	"github.com/victoryus84/gorders/internal/dto"
	"github.com/victoryus84/gorders/internal/models"
)

// 1. INTERFAȚA (Contractul) - Verifică litera 's' în Repository!
type ClientRepository interface {
	CreateClient(client *models.Client) error
	FindClientByFiscalID(fiscalID string) (*models.Client, error)
	GetFirst1000Clients() ([]models.Client, error)
	FindClientsByQuery(query string) ([]models.Client, error)
	FindClientByID(id uint) (*models.Client, error)
	CreateClientAddress(addr *models.ClientAddress) error
}

// Aici pui toate metodele pe care vrei să le folosească Handler-ul
type ClientService interface {
	ProcessClientImport(requests []dto.ClientDTO) dto.ImportResult
	GetFirst1000Clients() ([]models.Client, error)
	SearchClients(query string) ([]dto.ClientDTO, error)
	FindClientByID(id uint) (*models.Client, error)
	ProcessAddressImport(requests []dto.ClientAddressDTO, ownerID uint) dto.ImportResult
}

// 2. Facem structura PRIVATĂ (schimbăm 'C' mare în 'c' mic)
type clientService struct {
	repo ClientRepository
}

func NewClientService(repo ClientRepository) ClientService {
	return &clientService{repo: repo}
}

// ProcessClientImport - Logica masivă de import pe care am scos-o din Handler
func (s *clientService) ProcessClientImport(requests []dto.ClientDTO) dto.ImportResult {
	created := make([]*models.Client, 0)
	skipped := make([]map[string]string, 0)

	for _, req := range requests {
		// A. Validare de bază (Logica ta din API-ul vechi)
		if req.ClientTypeID == 0 || strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.FiscalID) == "" {
			skipped = append(skipped, map[string]string{"fiscal_id": req.FiscalID, "reason": "missing_required_fields"})
			continue
		}

		// B. Verificare duplicate
		existing, err := s.repo.FindClientByFiscalID(req.FiscalID)
		if err == nil && existing != nil {
			skipped = append(skipped, map[string]string{"fiscal_id": req.FiscalID, "reason": "duplicate"})
			continue
		}

		// C. Sanitizarea email-ului (Logica ta deșteaptă)
		emailPtr := s.sanitizeEmail(req.Email)

		// D. Mapare DTO -> Model
		client := &models.Client{
			ClientTypeID:  req.ClientTypeID,
			Name:          req.Name,
			FiscalID:      req.FiscalID,
			Email:         emailPtr,
			Phone:         req.Phone,
			FiscalAddress: req.FiscalAddress,
			PostalAddress: req.PostalAddress,
		}

		// E. Salvare
		if err := s.repo.CreateClient(client); err != nil {
			skipped = append(skipped, map[string]string{"fiscal_id": req.FiscalID, "reason": err.Error()})
			continue
		}
		created = append(created, client)
	}

	return dto.ImportResult{
		Status:        "success",
		TotalCreated:  len(created),
		TotalSkipped:  len(skipped),
		ErrorsPreview: s.limitErrors(skipped, 20),
		Message:       "Import finalizat",
	}
}

// SearchClients - Caută și mapează rezultatele în DTO-uri curate
func (s *clientService) SearchClients(query string) ([]dto.ClientDTO, error) {
	dbClients, err := s.repo.FindClientsByQuery(query)
	if err != nil {
		return nil, err
	}

	response := make([]dto.ClientDTO, len(dbClients))
	for i, cl := range dbClients {
		var emailStr string
		if cl.Email != nil {
			emailStr = *cl.Email
		}
		response[i] = dto.ClientDTO{
			ID:            cl.ID,
			ClientTypeID:  cl.ClientTypeID,
			Name:          cl.Name,
			FiscalID:      cl.FiscalID,
			Email:         emailStr,
			Phone:         cl.Phone,
			FiscalAddress: cl.FiscalAddress,
			PostalAddress: cl.PostalAddress,
		}
	}
	return response, nil
}

// ProcessAddressImport - Importul de adrese
func (s *clientService) ProcessAddressImport(requests []dto.ClientAddressDTO, ownerID uint) dto.ImportResult {
	createdCount := 0
	skipped := make([]map[string]string, 0)

	for _, req := range requests {
		if strings.TrimSpace(req.FiscalID) == "" {
			skipped = append(skipped, map[string]string{"address": req.Address, "reason": "fiscal_id_empty"})
			continue
		}

		client, err := s.repo.FindClientByFiscalID(req.FiscalID)
		if err != nil {
			skipped = append(skipped, map[string]string{"address": req.Address, "reason": "client_not_found"})
			continue
		}

		addr := &models.ClientAddress{
			ClientID: client.ID,
			Name:     req.Name,
			Address:  &req.Address,
			Type:     req.Type,
			OwnerID:  ownerID,
		}

		if err := s.repo.CreateClientAddress(addr); err != nil {
			skipped = append(skipped, map[string]string{"address": req.Address, "reason": err.Error()})
			continue
		}
		createdCount++
	}

	return dto.ImportResult{
		Status:        "success",
		TotalCreated:  createdCount,
		TotalSkipped:  len(skipped),
		ErrorsPreview: s.limitErrors(skipped, 20),
	}
}

// Metodele standard (Passthrough către repo)
func (s *clientService) GetFirst1000Clients() ([]models.Client, error) { return s.repo.GetFirst1000Clients() }

func (s *clientService) FindClientByID(id uint) (*models.Client, error) { return s.repo.FindClientByID(id) }

// --- Helpers Private ---

func (s *clientService) sanitizeEmail(email string) *string {
	raw := strings.TrimSpace(email)
	low := strings.ToLower(raw)
	if low == "" || low == "not inserted" || low == "n/a" || low == "none" {
		return nil
	}
	return &raw
}

func (s *clientService) limitErrors(skipped []map[string]string, limit int) []map[string]string {
	if len(skipped) > limit {
		return skipped[:limit]
	}
	return skipped
}