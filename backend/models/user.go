package models

import (
	"time"
)

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"uniqueIndex;not null" json:"username"`
	PasswordHash string    `gorm:"not null" json:"-"`
	Role         string    `gorm:"default:user" json:"role"` // 'admin' or 'user'
	CreatedAt    time.Time `json:"created_at"`
	
	// Contact Info (used for domain WHOIS data)
	ContactName    string `gorm:"default:''" json:"contact_name"`
	ContactOrg     string `gorm:"default:''" json:"contact_org"`
	ContactEmail   string `gorm:"default:''" json:"contact_email"`
	ContactPhone   string `gorm:"default:''" json:"contact_phone"`
	ContactAddress string `gorm:"default:''" json:"contact_address"`
	ContactCity    string `gorm:"default:''" json:"contact_city"`
	ContactState   string `gorm:"default:''" json:"contact_state"`
	ContactZip     string `gorm:"default:''" json:"contact_zip"`
	ContactCountry string `gorm:"default:''" json:"contact_country"`
}
