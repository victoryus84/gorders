package repository

import (
	"strings"

	"github.com/victoryus84/gorders/internal/models"
    "gorm.io/gorm/clause"
)

// Client methods
func (rep *Repository) CreateClient(client *models.Client) error {
    if client.Email != nil {
        em := strings.TrimSpace(*client.Email)
        el := strings.ToLower(em)

        // Dacă e mizerie (n/a, none, gol), îl facem NIL
        if em == "" || el == "not inserted" || el == "n/a" || el == "none" {
            client.Email = nil // În DB se va duce NULL (și e unic!)
        }
    }

    return rep.db.Create(client).Error
}

func (rep *Repository) UpsertClient(client *models.Client) error {
   // 1. Definim separat regula: "Ce facem dacă găsim un duplicat?"
    regulaConflict := clause.OnConflict{
        Columns: []clause.Column{{Name: "code"}}, // Cheia unică după care căutăm
        DoUpdates: clause.AssignmentColumns([]string{
            "name", "client_type_id", "fiscal_id", "email", "phone", 
            "fiscal_address", "postal_address", "client_group_id",
        }), // Câmpurile pe care le actualizăm
    }

    // 2. Aplicăm regula și executăm comanda de Salvare/Creare
    rezultat := rep.db.Clauses(regulaConflict).Create(client)

    // 3. Verificăm eroarea în stilul clasic
    if rezultat.Error != nil {
        return rezultat.Error
    }

    // 4. Totul a decurs perfect
    return nil
}

func (rep *Repository) UpsertClientsBatch(clients []*models.Client, batchSize int) error {
    regulaConflict := clause.OnConflict{
        Columns:   []clause.Column{{Name: "code"}},
        DoUpdates: clause.AssignmentColumns([]string{
            "name", "client_type_id", "fiscal_id", "email", "phone", 
            "fiscal_address", "postal_address", "client_group_id",
        }),
    }

    // Aici e magia: CreateInBatches
    return rep.db.Clauses(regulaConflict).CreateInBatches(clients, batchSize).Error
}

func (rep *Repository) CreateClientGroup(group *models.ClientGroup) error {
    return rep.db.Create(group).Error
}

func (rep *Repository) UpsertClientGroup(group *models.ClientGroup) error {
    // 1. Definim regula de conflict (ce facem dacă găsim un duplicat pe 'code')
    regulaConflict := clause.OnConflict{
        Columns: []clause.Column{{Name: "code"}}, // Cheia unică după care căutăm
        DoUpdates: clause.AssignmentColumns([]string{
            "name", // Câmpurile pe care le actualizăm dacă există deja
        }),
    }

    // 2. Aplicăm regula și executăm comanda de Salvare/Creare
    rezultat := rep.db.Clauses(regulaConflict).Create(group)

    // 3. Verificăm eroarea în stilul clasic
    if rezultat.Error != nil {
        return rezultat.Error
    }

    // 4. Totul a decurs perfect
    return nil
}

func (rep *Repository) CreateClientAddress(addr *models.ClientAddress) error {
	return rep.db.Create(addr).Error
}

func (rep *Repository) FindClientByCode(code string) (*models.Client, error) {
	var client models.Client
	err := rep.db.Where("code = ?", code).First(&client).Error
	return &client, err
}

func (rep *Repository) FindClientByFiscalID(fiscalID string) (*models.Client, error) {
	var client models.Client
	err := rep.db.Where("fiscal_id = ?", fiscalID).First(&client).Error
	return &client, err
}

func (rep *Repository) GetFirst1000Clients() ([]models.Client, error) {
	var clients []models.Client
	err := rep.db.Preload("ClientType").Limit(1000).Find(&clients).Error
	return clients, err
}

func (rep *Repository) FindClientsByQuery(query string) ([]models.Client, error) {
	if len(query) < 3 {
		return []models.Client{}, nil // Return empty if less than 3 chars
	}
	var clients []models.Client
	err := rep.db.
		Where("name ILIKE ? OR email ILIKE ? OR fiscal_id ILIKE ? OR phone ILIKE ?",
        "%"+query+"%",
        "%"+query+"%",
        "%"+query+"%",
        "%"+query+"%").
    Limit(50).
    Find(&clients).Error
	return clients, err
}

func (rep *Repository) FindClientByID(id uint) (*models.Client, error) {
	var client models.Client
	// Încercăm să găsim clientul și să încărcăm și tipul lui (Preload)
    err := rep.db.Preload("ClientType").First(&client, id).Error
    
    // 1. Dacă a apărut o eroare (nu există în DB sau e picat serverul)
    if err != nil {
        return nil, err // Returnăm "mâna goală" (nil) și eroarea
    }

    // 2. Dacă totul e ok, returnăm adresa obiectului plin
    return &client, nil
}

func (rep *Repository) FindClientGroupByName(name string) (*models.ClientGroup, error) {
    var group models.ClientGroup
    result := rep.db.Where("name = ?", name).First(&group)
    
    if result.Error != nil {
        return nil, result.Error
    }
    return &group, nil
}

func (rep *Repository) FindClientGroupByCode(code string) (*models.ClientGroup, error) {
    var group models.ClientGroup
    result := rep.db.Where("code = ?", code).First(&group)
    
    if result.Error != nil {
        return nil, result.Error
    }
    return &group, nil
}

func (rep *Repository) GetAllClientGroups() ([]models.ClientGroup, error) {
    var groups []models.ClientGroup
    if err := rep.db.Find(&groups).Error; err != nil {
        return nil, err
    }
    
    return groups, nil
}








