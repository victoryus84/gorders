package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/victoryus84/gorders/internal/config"
	"github.com/victoryus84/gorders/internal/dto"
	"github.com/victoryus84/gorders/internal/kafka"
	"github.com/victoryus84/gorders/internal/logger"
	"github.com/victoryus84/gorders/internal/models"
)

// Client repository-related database operations
type ClientRepository interface {
	// Client methods
	CreateClient(client *models.Client) error
	UpsertClient(client *models.Client) error
	UpsertClientsBatch(clients []*models.Client, batchSize int) error
	FindClientByCode(code string) (*models.Client, error)
	FindClientByFiscalID(fiscalID string) (*models.Client, error)
	GetFirst1000Clients() ([]models.Client, error)
	FindClientsByQuery(query string) ([]models.Client, error)
	FindClientByID(id uint) (*models.Client, error)
	FindAllClientCodesMap() (map[string]uint, error)
	// Group client methods
	CreateClientGroup(group *models.ClientGroup) error
	UpsertClientGroup(group *models.ClientGroup) error
	FindClientGroupByName(name string) (*models.ClientGroup, error)
	FindClientGroupByCode(code string) (*models.ClientGroup, error)
	GetAllClientGroups() ([]models.ClientGroup, error)
	// Address client methods
	CreateClientAddress(addr *models.ClientAddress) error
	UpsertClientsAddressBatch(addr []*models.ClientAddress, batchSize int) error
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
	// 1. Pregătim "cutia" în care adunăm clienții
	clientsToSave := make([]*models.Client, 0, len(requests))
	skipped := make([]map[string]string, 0)
	topic := svc.cfg.GetTopic("clients")

	// Extragem toate grupele din baza de date pentru a le avea la îndemână
	groupMap := make(map[string]uint)
	if existingGroups, err := svc.rep.GetAllClientGroups(); err == nil {
		for _, g := range existingGroups {
			groupMap[g.Code] = g.ID
		}
	}

	for _, req := range requests {
		// A. Validare de bază (Logica ta din API-ul vechi)
		if req.ClientTypeID == 0 || strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.FiscalID) == "" {
			skipped = append(skipped, map[string]string{"fiscal_id": req.FiscalID, "reason": "missing_required_fields"})
			continue
		}

		// B. Verificare duplicate ()
		// existing, err := svc.rep.FindClientByCode(req.Code)
		// if err == nil && existing != nil {
		// 	skipped = append(skipped, map[string]string{"fiscal_id": req.FiscalID, "reason": "duplicate"})
		// 	continue
		// }

		// C. Sanitizarea cimpurilor (ex: email) și maparea codului grupei în ID-ul real din DB
		emailPtr := svc.sanitizeEmail(req.Email)
		dbGroupID := svc.resolveGroupIDFromMap(req.GroupCode, groupMap)

		// D. Mapare DTO -> Model
		client := &models.Client{
			Code:          req.Code,
			ClientTypeID:  req.ClientTypeID,
			Name:          req.Name,
			FiscalID:      req.FiscalID,
			Email:         emailPtr,
			Phone:         req.Phone,
			FiscalAddress: req.FiscalAddress,
			PostalAddress: req.PostalAddress,
			ClientGroupID: dbGroupID,
		}

		clientsToSave = append(clientsToSave, client)
	}

	// 3. Executăm UPSERT-ul masiv (câte 1000 o dată)
	// Durează o fracțiune de secundă pentru toți cei X de clienți!
	if err := svc.rep.UpsertClientsBatch(clientsToSave, 1000); err != nil {
		return dto.ImportResult{
			Status:  "error",
			Message: "Eroare fatală la salvarea în masă: " + err.Error(),
		}
	}

	// 4. KAFKA: Acum că știm că s-au salvat 100% cu succes în Postgres, îi trimitem către Kafka
	for _, mod_client := range clientsToSave {
		go func(client *models.Client) {
			payload, _ := json.Marshal(client)
			_ = svc.kfk.Publish(context.Background(), topic, client.Code, payload)
		}(mod_client)
	}

	return dto.ImportResult{
		Status:         "success",
		TotalProcessed: len(clientsToSave),
		TotalSkipped:   len(skipped),
		ErrorsPreview:  svc.limitErrors(skipped, 20),
		Message:        "Syncronizare clienti finalizata!",
	}
}

func (svc *clientService) ProcessClientGroupImport(requests []dto.ClientGroupDTO) dto.ImportResult {
	created := make([]*models.ClientGroup, 0)
	skipped := make([]map[string]string, 0)
	topic := svc.cfg.GetTopic("client_groups")

	// 1. INIȚIALIZAREA: Încărcăm grupele existente o singură dată
	groupMap := make(map[string]uint)
	if existingGroups, err := svc.rep.GetAllClientGroups(); err == nil {
		for _, g := range existingGroups {
			groupMap[g.Code] = g.ID
		}
	}

	for _, req := range requests {
		// A. Validare de bază
		if strings.TrimSpace(req.Name) == "" {
			skipped = append(skipped, map[string]string{"name": req.Name, "reason": "missing_required_fields"})
			continue
		}

		// B. Verificare duplicate
		// existing, err := svc.rep.FindClientGroupByCode(req.Code)
		// if err == nil && existing != nil {
		//     skipped = append(skipped, map[string]string{"name": req.Name, "reason": "duplicate"})
		//     continue
		// }

		// C. MAPAREA IERARHIEI - Căutăm în dicționar dacă avem codul părintelui și luăm ID-ul lui
		var parentIDPtr *uint
		if cleanCode := strings.TrimSpace(req.ParentCode); cleanCode != "" && cleanCode != "not inserted" {
			if id, exists := groupMap[cleanCode]; exists {
				parentIDPtr = &id
			}
		}

		// D. Mapare DTO -> Model
		clientgroup := &models.ClientGroup{
			Code:        req.Code,
			Name:        req.Name,
			Description: req.Description,
			ParentID:    parentIDPtr, // Acum ia adresa reală sau rămâne nil
		}

		// E. Salvare
		if err := svc.rep.UpsertClientGroup(clientgroup); err != nil {
			skipped = append(skipped, map[string]string{"code": req.Code, "reason": "upsert_failed: " + err.Error()})
			continue
		}

		// F. ACTUALIZAREA: Scriem grupa abia salvată în dicționar pentru viitorii ei copii!
		groupMap[clientgroup.Code] = clientgroup.ID

		// G. KAFKA
		go func(mod *models.ClientGroup) {
			payload, _ := json.Marshal(mod)
			_ = svc.kfk.Publish(context.Background(), topic, mod.Name, payload)
		}(clientgroup)

		created = append(created, clientgroup)
	}

	return dto.ImportResult{
		Status:         "success",
		TotalProcessed: len(created),
		TotalSkipped:   len(skipped),
		ErrorsPreview:  svc.limitErrors(skipped, 20),
		Message:        "Import finalizat",
	}
}

// ProcessAddressImport - Importul de adrese
func (svc *clientService) ProcessAddressImport(requests []dto.ClientAddressDTO, ownerID uint) dto.ImportResult {
	skipped := make([]map[string]string, 0)

	// ==========================================
	// 1. DICȚIONARUL ÎN RAM (1 singur drum la DB!)
	// ==========================================
	clientMap, err := svc.rep.FindAllClientCodesMap()
	if err != nil {
		clientMap = make(map[string]uint)
	}

	// AICI E "SITA": Pregătim un map în loc de array pentru a elimina dublurile automat
	uniqueAddresses := make(map[string]*models.ClientAddress)

	// ==========================================
	// 2. PROCESĂM ÎN MEMORIE ȘI FILTRĂM DUBURILE
	// ==========================================
	for _, req := range requests {

		if strings.TrimSpace(req.Code) == "" {
			skipped = append(skipped, logSkip(req.Name, "Code (cod client) lipsă"))
			continue
		}

		dbClientID, exists := clientMap[req.Code]
		if !exists {
			skipped = append(skipped, logSkip(req.Name, "Adresa inexistenta: "+req.Code))
			continue
		}

		// Conversie DTO -> Model (asigură-te că în mapDTOToModel asignezi contract.SyncID = req.SyncID)
		address := mapDTOToModel(req, dbClientID, ownerID)

		// Get the key directly from 1C
		syncKey := req.SyncID

		// Check if 1C sent a duplicate in the same file
		if _, exists := uniqueAddresses[syncKey]; exists {
			// Log in English for the server console
			logger.LogInfo(fmt.Sprintf("🔥 XML DUPLICATE: Contract with SyncID '%s' appeared multiple times in the payload!", syncKey))

			// Textul pentru 1C rămâne în română (fără diacritice) ca să-l înțeleagă operatorii
			skipped = append(skipped, logSkip(req.Name, "(Addresses) Dublura in XML: "+syncKey))
		}

		// Add to the map
		uniqueAddresses[syncKey] = address
	}
	// --- RECONSTRUIM ARRAY-UL PENTRU GORM ---
	// Acum că am scăpat de dubluri, mutăm datele din map înapoi într-un slice
	addressesToSave := make([]*models.ClientAddress, 0, len(uniqueAddresses))
	for _, c := range uniqueAddresses {
		addressesToSave = append(addressesToSave, c)
	}

	// ==========================================
	// 3. SALVAREA ÎN MASĂ (Tunul!)
	// ==========================================
	if len(addressesToSave) > 0 {
		if err := svc.rep.UpsertClientsAddressBatch(addressesToSave, 1000); err != nil {
			return dto.ImportResult{
				Status:  "error",
				Message: "Eroare DB la salvarea în masă: " + err.Error(),
			}
		}
	}

	return dto.ImportResult{
		Status:         "success",
		TotalProcessed: len(addressesToSave), // Acum arată nr real, fără dubluri!
		TotalSkipped:   len(skipped),
		ErrorsPreview:  svc.limitErrors(skipped, 20),
		Message:        "Sincronizare adrese finalizata!",
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
	// 1. Curățăm codul (trim, lower) pentru a evita problemele de formatare și caze sensitivity
	cleanCode := strings.TrimSpace(code)
	lowerCode := strings.ToLower(cleanCode)

	// 2. Dacă nu a venit niciun cod, returnăm nil direct
	switch lowerCode {
	case "", "none", "n/a", "not inserted":
		return nil
	}

	// 3. AICI FOLOSIM PARAMETRUL: Căutăm codul în dicționar
	if id, exists := groupMap[cleanCode]; exists {
		return &id // Am găsit grupa, îi returnăm ID-ul sub formă de pointer
	}

	// 4. Dacă codul există în 1C dar nu l-am găsit în dicționarul nostru, returnăm nil
	return nil

}

// --- HELPER FUNCTIONS (Curățenia din labirint) ---

func mapDTOToModel(req dto.ClientAddressDTO, clientID uint, ownerID uint) *models.ClientAddress {
	return &models.ClientAddress{
		SyncID:       req.SyncID,
		Name:         req.Name,
		Address:      stringPtr(req.Address), // Aici folosim unealta universala!
		Type:         req.Type,
		DeliveryDays: req.DeliveryDays, // Pura magie: direct din 1C
		ClientID:     clientID,
		OwnerID:      ownerID,
	}
}
