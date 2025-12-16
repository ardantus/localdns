package models

import (
	"time"
)

type Domain struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"uniqueIndex;not null" json:"name"`
	UserID    uint      `gorm:"not null" json:"user_id"`
    User      User      `json:"user,omitempty"` // Association
	CreatedAt time.Time `json:"created_at"`
    // Relations
    Records   []Record  `json:"records,omitempty"`
}

type Record struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	DomainID  uint      `gorm:"not null;index" json:"domain_id"`
	Name      string    `gorm:"not null" json:"name"`
	Type      string    `gorm:"not null" json:"type"` // A, CNAME, etc.
	Content   string    `gorm:"not null" json:"content"`
	TTL       int       `gorm:"default:360" json:"ttl"`
	Prio      int       `gorm:"default:0" json:"prio"`
	Disabled  bool      `gorm:"default:false" json:"disabled"`
	CreatedAt time.Time `json:"created_at"`
}
