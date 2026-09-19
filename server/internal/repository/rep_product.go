package repository

import (
	"github.com/victoryus84/gorders/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause" // <-- MAGIC: Obligatoriu pentru salvarea masivă (Upsert)
)

// INTERFAȚA: Asta este exact ce așteaptă ProductService
type ProductRepository interface {
	CreateProduct(product *models.Product) error
	FindProductGroupByID(id uint) (*models.ProductGroup, error)
	FindProductByID(id uint) (*models.Product, error)
	FindProductsByQuery(query string) ([]models.Product, error)
	// FUNCȚIA NOUĂ PENTRU 1C: Salvare masivă!
	UpsertProductsBatch(items []models.Product, batchSize int) error
	GetFirst1000Products() ([]models.Product, error) 
}

// STRUCTURA
type productRepository struct {
	db *gorm.DB
}

// CONSTRUCTORUL (Observă că returnează interfața ProductRepository)
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

// ==========================================
// METODELE
// ==========================================

func (rep *productRepository) CreateProduct(product *models.Product) error {
	return rep.db.Create(product).Error
}

func (rep *productRepository) FindProductGroupByID(id uint) (*models.ProductGroup, error) {
	var group models.ProductGroup
	err := rep.db.First(&group, id).Error
	return &group, err
}

func (rep *productRepository) FindProductByID(id uint) (*models.Product, error) {
	var product models.Product
	err := rep.db.First(&product, id).Error
	return &product, err
}

// [CORECȚIE AICI]: Am schimbat coloanele din ILIKE! Căutăm după Nume, Cod sau Articol.
func (rep *productRepository) FindProductsByQuery(query string) ([]models.Product, error) {
	if len(query) < 3 {
		return []models.Product{}, nil
	}
	var products []models.Product
	
	// Acum căutăm doar prin câmpurile reale ale produsului!
	err := rep.db.
		Where("name ILIKE ? OR code ILIKE ? OR article ILIKE ?",
			"%"+query+"%", "%"+query+"%", "%"+query+"%").
		Limit(50).Find(&products).Error
	return products, err
}

// [NOU]: Motorul de sincronizare din 1C!
func (rep *productRepository) UpsertProductsBatch(items []models.Product, batchSize int) error {
	if len(items) == 0 {
		return nil
	}

	// clause.OnConflict spune bazei de date: "Dacă există deja acest cod (Code), fă-i update la restul datelor!"
	return rep.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}}, // Asigură-te că tabela ta are constrângere de unicitate pe coloana 'code'
		UpdateAll: true,
	}).CreateInBatches(items, batchSize).Error
}

func (rep *productRepository) GetFirst1000Products() ([]models.Product, error) {
	var products []models.Product
	err := rep.db.Limit(1000).Find(&products).Error
	return products, err
}