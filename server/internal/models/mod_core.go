package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UUIDModel este o structură de bază pe care o vom "îngropa" în alte modele.
type UUIDModel struct {
	UUID string `gorm:"type:uuid;uniqueIndex;default:null"`
}

// Punem și hook-ul aici, pentru că aparține de UUIDModel.
// Orice alt model care conține UUIDModel va rula automat acest cod.
func (m *UUIDModel) BeforeCreate(tx *gorm.DB) (err error) {
	if m.UUID == "" {
		m.UUID = uuid.New().String()
	}
	return nil
}