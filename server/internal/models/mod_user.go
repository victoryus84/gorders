package models

import (
	"gorm.io/gorm"
)

// ********** User - Utilizatorul sistemului **********
type User struct {
	gorm.Model
	UUIDModel `gorm:"embedded"`
	Email     string    `gorm:"unique;not null"`           // Email-ul utilizatorului (unic)
	Password  string    `gorm:"not null"`                  // Hash-ul parolei (nu se afișează în JSON)
	Role      string    `gorm:"type:varchar(20);not null"` // Rolul utilizatorului ("admin", "user" etc.)
	Channels  []Channel `gorm:"many2many:user_channels;"`  // Canalele de vânzări la care are acces utilizatorul
}

// ****************************************************

// AfterCreate hook pentru User - dacă este primul utilizator creat, îi setează rolul de admin
func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	if u.ID == 1 {
		tx.Model(u).Update("role", "admin")
	}
	return
}
