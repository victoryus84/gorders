package service

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/victoryus84/gorders/internal/config"
	"github.com/victoryus84/gorders/internal/dto"
	"github.com/victoryus84/gorders/internal/kafka"
	"github.com/victoryus84/gorders/internal/models"
)

// Client repository-related database operations
type ClientRepository interface {
	// Client methods
	CreateClient(client *models.Client) error
	FindClientByFiscalID(fiscalID string) (*models.Client, error)
	GetFirst1000Clients() ([]models.Client, error)
	FindClientsByQuery(query string) ([]models.Client, error)
	FindClientByID(id uint) (*models.Client, error)
	// Group client methods
	CreateClientGroup(group *models.ClientGroup) error
	FindClientGroupByName(name string) (*models.ClientGroup, error)
	FindClientGroupByCode(code string) (*models.ClientGroup, error)
	GetAllClientGroups() ([]models.ClientGroup, error)
	// Address client methods
	CreateClientAddress(addr *models.ClientAddress) error
}

// Aici pui toate metodele pe care vrei să le folosească Handler-ul
type ClientService interface {
	ProcessClientImport(requests []dto.ClientDTO) dto.ImportResult
	ProcessClientGroupImport(requests []dto.ClientGroupDTO) dto.ImportResult
	ProcessAddressImport(requests []dto.ClientAddressDTO, ownerID uint) dto.ImportResult
	SvcGetFirst1000Clients() ([]models.Client, error)
	SearchClients(query string) ([]dto.ClientDTO, error)
	SearchClientByID(id uint) (*models.Client, error)
}

// 2. Facem structura PRIVATĂ (schimbăm 'C' mare în 'c' mic)
type clientService struct {
	rep ClientRepository
	cfg *config.Config
	kfk *kafka.Producer
}

func NewClientService(
	rep ClientRepository,
	cfg *config.Config,
	kfk *kafka.Producer) ClientService {
	return &clientService{
		rep: rep,
		cfg: cfg,
		kfk: kfk}
}

// ProcessClientImport - Logica masivă de import pe care am scos-o din Handler
func (svc *clientService) ProcessClientImport(requests []dto.ClientDTO) dto.ImportResult {
	created := make([]*models.Client, 0)
	skipped := make([]map[string]string, 0)
	topic := svc.cfg.GetTopic("clients")

	for _, req := range requests {
		// A. Validare de bază (Logica ta din API-ul vechi)
		if req.ClientTypeID == 0 || strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.FiscalID) == "" {
			skipped = append(skipped, map[string]string{"fiscal_id": req.FiscalID, "reason": "missing_required_fields"})
			continue
		}

		// B. Verificare duplicate
		existing, err := svc.rep.FindClientByFiscalID(req.FiscalID)
		if err == nil && existing != nil {
			skipped = append(skipped, map[string]string{"fiscal_id": req.FiscalID, "reason": "duplicate"})
			continue
		}

		// C. Sanitizarea email-ului (Logica ta deșteaptă)
		emailPtr := svc.sanitizeEmail(req.Email)
		dbGroupID := svc.resolveGroupIDFromMap(req.GroupCode, nil) 

		// D. Mapare DTO -> Model
		client := &models.Client{
			ClientTypeID:  req.ClientTypeID,
			Name:          req.Name,
			FiscalID:      req.FiscalID,
			Email:         emailPtr,
			Phone:         req.Phone,
			FiscalAddress: req.FiscalAddress,
			PostalAddress: req.PostalAddress,
			ClientGroupID: dbGroupID,
		}

		// E. Salvare
		if err := svc.rep.CreateClient(client); err != nil {
			skipped = append(skipped, map[string]string{"fiscal_id": req.FiscalID, "reason": err.Error()})
			continue
		}

		// F. KAFKA
		// Folosim gorutină pentru a nu încetini importul
		go func(mod_client *models.Client) {
			payload, _ := json.Marshal(mod_client)

			// Publish-ul tău universal
			_ = svc.kfk.Publish(context.Background(), topic, mod_client.FiscalID, payload)

			// Notă: Dacă ai un logger (ex: zap, logrus), aici ar trebui să fie:
			// svc.logger.Error("kafka publish failed", "client", mod_client.FiscalID, "err", err)
		}(client)

		created = append(created, client)
	}

	return dto.ImportResult{
		Status:        "success",
		TotalCreated:  len(created),
		TotalSkipped:  len(skipped),
		ErrorsPreview: svc.limitErrors(skipped, 20),
		Message:       "Import finalizat",
	}
}

func (svc *clientService) ProcessClientGroupImport(requests []dto.ClientGroupDTO) dto.ImportResult {
	created := make([]*models.ClientGroup, 0)
	skipped := make([]map[string]string, 0)
	topic := svc.cfg.GetTopic("client_groups")

	for _, req := range requests {
		// A. Validare de bază (Logica ta din API-ul vechi)
		if strings.TrimSpace(req.Name) == "" {
			skipped = append(skipped, map[string]string{"name": req.Name, "reason": "missing_required_fields"})
			continue
		}

		// B. Verificare duplicate
		existing, err := svc.rep.FindClientGroupByCode(req.Code)
		if err == nil && existing != nil {
			skipped = append(skipped, map[string]string{"name": req.Name, "reason": "duplicate"})
			continue
		}

		// D. Mapare DTO -> Model
		clientgroup := &models.ClientGroup{
			Name:        req.Name,
			Description: req.Description,
		}

		// E. Salvare
		if err := svc.rep.CreateClientGroup(clientgroup); err != nil {
			skipped = append(skipped, map[string]string{"name": req.Name, "reason": err.Error()})
			continue
		}

		// F. KAFKA
		// Folosim gorutină pentru a nu încetini importul
		go func(mod *models.ClientGroup) {
			payload, _ := json.Marshal(mod)

			// Publish-ul tău universal
			_ = svc.kfk.Publish(context.Background(), topic, mod.Name, payload)

			// Notă: Dacă ai un logger (ex: zap, logrus), aici ar trebui să fie:
			// svc.logger.Error("kafka publish failed", "client", mod.FiscalID, "err", err)
		}(clientgroup)

		created = append(created, clientgroup)
	}

	return dto.ImportResult{
		Status:        "success",
		TotalCreated:  len(created),
		TotalSkipped:  len(skipped),
		ErrorsPreview: svc.limitErrors(skipped, 20),
		Message:       "Import finalizat",
	}
}

// ProcessAddressImport - Importul de adrese
func (svc *clientService) ProcessAddressImport(requests []dto.ClientAddressDTO, ownerID uint) dto.ImportResult {
	createdCount := 0
	skipped := make([]map[string]string, 0)

	for _, req := range requests {
		if strings.TrimSpace(req.FiscalID) == "" {
			skipped = append(skipped, map[string]string{"address": req.Address, "reason": "fiscal_id_empty"})
			continue
		}

		client, err := svc.rep.FindClientByFiscalID(req.FiscalID)
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

		if err := svc.rep.CreateClientAddress(addr); err != nil {
			skipped = append(skipped, map[string]string{"address": req.Address, "reason": err.Error()})
			continue
		}
		createdCount++
	}

	return dto.ImportResult{
		Status:        "success",
		TotalCreated:  createdCount,
		TotalSkipped:  len(skipped),
		ErrorsPreview: svc.limitErrors(skipped, 20),
	}
}

// SearchClients - Caută și mapează rezultatele în DTO-uri curate
func (svc *clientService) SearchClients(query string) ([]dto.ClientDTO, error) {
	dbClients, err := svc.rep.FindClientsByQuery(query)
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

// Metodele standard (Passthrough către repo)
func (svc *clientService) SvcGetFirst1000Clients() ([]models.Client, error) {
	return svc.rep.GetFirst1000Clients()
}

func (svc *clientService) SearchClientByID(id uint) (*models.Client, error) {
	return svc.rep.FindClientByID(id)
}

// --- Helpers Private ---

func (svc *clientService) sanitizeEmail(email string) *string {
	raw := strings.TrimSpace(email)
	low := strings.ToLower(raw)
	if low == "" || low == "not inserted" || low == "n/a" || low == "none" {
		return nil
	}
	return &raw
}

func (svc *clientService) limitErrors(skipped []map[string]string, limit int) []map[string]string {
	if len(skipped) > limit {
		return skipped[:limit]
	}
	return skipped
}

func (svc *clientService) resolveGroupIDFromMap(code string, groupMap map[string]uint) *uint {
	cleanCode := strings.TrimSpace(code)

	if cleanCode == "" || strings.ToLower(cleanCode) == "none" || strings.ToLower(cleanCode) == "n/a" {
		return nil
	}

	group, err := svc.rep.FindClientGroupByCode(cleanCode)
	if err != nil || group == nil {
		return nil
	}

	return &group.ID
}