package repository

import (
	"strings"

	"github.com/victoryus84/gorders/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 1. INTERFAȚA (Mutată aici din service)
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
	GetAllClientGroups() ([]models.ClientGroup, error)
	// Address client methods
	CreateClientAddress(addr *models.ClientAddress) error
	UpsertClientsAddressBatch(addr []*models.ClientAddress, batchSize int) error
}

// 2. STRUCTURA PRIVATĂ
type clientRepository struct { // presupunând că așa se numește
	db *gorm.DB
}

// 3. CONSTRUCTORUL
// Asigură-te că returnează interfața ClientRepository, nu pointer la structură
func NewClientRepository(db *gorm.DB) ClientRepository {
	return &clientRepository{db: db}
}

// ==========================================
// METODE SPECIFICE PENTRU CLIENȚI
// ==========================================

func (rep *clientRepository) CreateClient(client *models.Client) error {
	if client.Email != nil {
		em := strings.TrimSpace(*client.Email)
		el := strings.ToLower(em)

		if em == "" || el == "not inserted" || el == "n/a" || el == "none" {
			client.Email = nil
		}
	}
	return rep.db.Create(client).Error
}

func (rep *clientRepository) UpsertClient(client *models.Client) error {
	return rep.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "code"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"name", "client_type_id", "fiscal_id", "email", "phone",
			"fiscal_address", "postal_address", "client_group_id",
		}),
	}).Create(client).Error
}

func (rep *clientRepository) UpsertClientsBatch(clients []*models.Client, batchSize int) error {
	return rep.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "code"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"name", "client_type_id", "fiscal_id", "email", "phone",
			"fiscal_address", "postal_address", "client_group_id",
		}),
	}).CreateInBatches(clients, batchSize).Error
}

func (rep *clientRepository) FindClientByID(id uint) (*models.Client, error) {
	var client models.Client
	err := rep.db.Preload("ClientType").First(&client, id).Error
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (rep *clientRepository) FindClientByCode(code string) (*models.Client, error) {
	var client models.Client
	// Căutăm primul client care are codul primit din 1C
	err := rep.db.Where("code = ?", code).First(&client).Error
	return &client, err
}

func (rep *clientRepository) FindClientByFiscalID(fiscalID string) (*models.Client, error) {
	var client models.Client
	// Căutăm clientul după IDNO / Codul Fiscal specific Republicii Moldova
	err := rep.db.Where("fiscal_id = ?", fiscalID).First(&client).Error
	return &client, err
}

func (rep *clientRepository) GetFirst1000Clients() ([]models.Client, error) {
	var clients []models.Client
	err := rep.db.Preload("ClientType").Limit(1000).Find(&clients).Error
	return clients, err
}

func (rep *clientRepository) FindClientsByQuery(query string) ([]models.Client, error) {
	if len(query) < 3 {
		return []models.Client{}, nil
	}
	var clients []models.Client
	err := rep.db.
		Where("name ILIKE ? OR email ILIKE ? OR fiscal_id ILIKE ? OR phone ILIKE ?",
			"%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%").
		Limit(50).Find(&clients).Error
	return clients, err
}

func (rep *clientRepository) FindAllClientCodesMap() (map[string]uint, error) {
	type clientResult struct {
		ID   uint
		Code string
	}
	var results []clientResult

	if err := rep.db.Model(&models.Client{}).Select("id, code").Find(&results).Error; err != nil {
		return nil, err
	}

	clientMap := make(map[string]uint, len(results))
	for _, r := range results {
		clientMap[r.Code] = r.ID
	}
	return clientMap, nil
}

// ==========================================
// METODE SPECIFICE PENTRU GRUPURI
// ==========================================
func (rep *clientRepository) CreateClientGroup(group *models.ClientGroup) error {
	return rep.db.Create(group).Error
}

func (rep *clientRepository) UpsertClientGroup(group *models.ClientGroup) error {
	return rep.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}},
		DoUpdates: clause.AssignmentColumns([]string{"name"}),
	}).Create(group).Error
}

func (rep *clientRepository) GetAllClientGroups() ([]models.ClientGroup, error) {
    var groups []models.ClientGroup
    err := rep.db.Find(&groups).Error
    return groups, err
}

// ==========================================
// METODE SPECIFICE PENTRU ADRESE
// ==========================================
func (rep *clientRepository) CreateClientAddress(addr *models.ClientAddress) error {
	return rep.db.Create(addr).Error
}

func (rep *clientRepository) UpsertClientsAddressBatch(addresses []*models.ClientAddress, batchSize int) error {
	return rep.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "sync_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"name", "client_id", "address", "delivery_days", "type", "updated_at",
		}),
	}).CreateInBatches(addresses, batchSize).Error
}