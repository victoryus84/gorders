package repository

import (
	"strings"		

	"github.com/victoryus84/gorders/internal/models"
)

// Client methods
func (repository *Repository) CreateClient(client *models.Client) error {
    if client.Email != nil {
        em := strings.TrimSpace(*client.Email)
        el := strings.ToLower(em)

        // Dacă e mizerie (n/a, none, gol), îl facem NIL
        if em == "" || el == "not inserted" || el == "n/a" || el == "none" {
            client.Email = nil // În DB se va duce NULL (și e unic!)
        }
    }

    return repository.db.Create(client).Error
}
func (repository *Repository) GetFirst1000Clients() ([]models.Client, error) {
	var clients []models.Client
	err := repository.db.Preload("ClientType").Limit(1000).Find(&clients).Error
	return clients, err
}
func (repository *Repository) FindClientsByQuery(query string) ([]models.Client, error) {
	if len(query) < 3 {
		return []models.Client{}, nil // Return empty if less than 3 chars
	}
	var clients []models.Client
	err := repository.db.
		Where("name ILIKE ? OR email ILIKE ? OR fiscal_id ILIKE ? OR phone ILIKE ?",
        "%"+query+"%",
        "%"+query+"%",
        "%"+query+"%",
        "%"+query+"%").
    Limit(50).
    Find(&clients).Error
	return clients, err
}
func (repository *Repository) FindClientByID(id uint) (*models.Client, error) {
	var client models.Client
	// Încercăm să găsim clientul și să încărcăm și tipul lui (Preload)
    err := repository.db.Preload("ClientType").First(&client, id).Error
    
    // 1. Dacă a apărut o eroare (nu există în DB sau e picat serverul)
    if err != nil {
        return nil, err // Returnăm "mâna goală" (nil) și eroarea
    }

    // 2. Dacă totul e ok, returnăm adresa obiectului plin
    return &client, nil
}
func (repository *Repository) FindClientByFiscalID(fiscalID string) (*models.Client, error) {
	var client models.Client
	err := repository.db.Where("fiscal_id = ?", fiscalID).First(&client).Error
	return &client, err
}
func (repository *Repository) CreateClientAddress(addr *models.ClientAddress) error {
	return repository.db.Create(addr).Error
}





