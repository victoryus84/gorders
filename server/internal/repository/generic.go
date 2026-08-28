package repository

import "gorm.io/gorm"

// GenericRepo este structura universală pentru orice model din DB
type GenericRepo[T any] struct {
	db *gorm.DB
}

// NewGenericRepo este constructorul
func NewGenericRepo[T any](db *gorm.DB) *GenericRepo[T] {
	return &GenericRepo[T]{db: db}
}

// Create salvează un obiect nou în baza de date
func (r *GenericRepo[T]) Create(entity *T) error {
	return r.db.Create(entity).Error
}

// FindByID caută un obiect după ID (presupunând că ai o coloană 'id')
func (r *GenericRepo[T]) FindByID(id string) (*T, error) {
	var entity T
	err := r.db.First(&entity, "id = ?", id).Error
	return &entity, err
}

// FindAll aduce toate înregistrările din tabelă
func (r *GenericRepo[T]) FindAll() ([]T, error) {
	var entities []T
	err := r.db.Find(&entities).Error
	return entities, err
}

// Update actualizează un obiect existent
func (r *GenericRepo[T]) Update(entity *T) error {
	return r.db.Save(entity).Error
}

// Delete șterge un obiect
func (r *GenericRepo[T]) Delete(entity *T) error {
	return r.db.Delete(entity).Error
}

// 🚀 BONUS de Senior: O funcție prin care poți rula query-uri custom direct din Service!
func (r *GenericRepo[T]) DB() *gorm.DB {
	return r.db
}
