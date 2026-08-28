package repository

import (
	"strings"

	"github.com/victoryus84/gorders/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ClientRepository struct {
	db *gorm.DB
}

// Constructorul pentru Uber Fx
func NewClientRepository(db *gorm.DB) *ClientRepository {
	return &ClientRepository{db: db}
}

// ==========================================
// METODE SPECIFICE PENTRU CLIENȚI
// ==========================================

func (rep *ClientRepository) CreateClient(client *models.Client) error {
	if client.Email != nil {
		em := strings.TrimSpace(*client.Email)
		el := strings.ToLower(em)

		if em == "" || el == "not inserted" || el == "n/a" || el == "none" {
			client.Email = nil
		}
	}
	return rep.db.Create(client).Error
}

func (rep *ClientRepository) UpsertClient(client *models.Client) error {
	return rep.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "code"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"name", "client_type_id", "fiscal_id", "email", "phone",
			"fiscal_address", "postal_address", "client_group_id",
		}),
	}).Create(client).Error
}

func (rep *ClientRepository) UpsertClientsBatch(clients []*models.Client, batchSize int) error {
	return rep.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "code"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"name", "client_type_id", "fiscal_id", "email", "phone",
			"fiscal_address", "postal_address", "client_group_id",
		}),
	}).CreateInBatches(clients, batchSize).Error
}

func (rep *ClientRepository) FindClientByID(id uint) (*models.Client, error) {
	var client models.Client
	err := rep.db.Preload("ClientType").First(&client, id).Error
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (rep *ClientRepository) GetFirst1000Clients() ([]models.Client, error) {
	var clients []models.Client
	err := rep.db.Preload("ClientType").Limit(1000).Find(&clients).Error
	return clients, err
}

func (rep *ClientRepository) FindClientsByQuery(query string) ([]models.Client, error) {
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

func (rep *ClientRepository) FindAllClientCodesMap() (map[string]uint, error) {
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

func (rep *ClientRepository) UpsertClientGroup(group *models.ClientGroup) error {
	return rep.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}},
		DoUpdates: clause.AssignmentColumns([]string{"name"}),
	}).Create(group).Error
}

// ==========================================
// METODE SPECIFICE PENTRU ADRESE
// ==========================================

func (rep *ClientRepository) UpsertClientsAddressBatch(addresses []*models.ClientAddress, batchSize int) error {
	return rep.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "sync_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"name", "client_id", "address", "delivery_days", "type", "updated_at",
		}),
	}).CreateInBatches(addresses, batchSize).Error
}