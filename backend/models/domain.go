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
	UpdatedAt time.Time `json:"updated_at"`
	ExpiresAt time.Time `json:"expires_at"`
	
	// Registrant Contact Info (WHOIS data)
	RegistrantName    string `gorm:"default:''" json:"registrant_name"`
	RegistrantOrg     string `gorm:"default:''" json:"registrant_org"`
	RegistrantEmail   string `gorm:"default:''" json:"registrant_email"`
	RegistrantPhone   string `gorm:"default:''" json:"registrant_phone"`
	RegistrantAddress string `gorm:"default:''" json:"registrant_address"`
	RegistrantCity    string `gorm:"default:''" json:"registrant_city"`
	RegistrantState   string `gorm:"default:''" json:"registrant_state"`
	RegistrantZip     string `gorm:"default:''" json:"registrant_zip"`
	RegistrantCountry string `gorm:"default:''" json:"registrant_country"`
	
	// Admin Contact Info (defaults to Registrant if empty)
	AdminName    string `gorm:"default:''" json:"admin_name"`
	AdminOrg     string `gorm:"default:''" json:"admin_org"`
	AdminEmail   string `gorm:"default:''" json:"admin_email"`
	AdminPhone   string `gorm:"default:''" json:"admin_phone"`
	AdminAddress string `gorm:"default:''" json:"admin_address"`
	AdminCity    string `gorm:"default:''" json:"admin_city"`
	AdminState   string `gorm:"default:''" json:"admin_state"`
	AdminZip     string `gorm:"default:''" json:"admin_zip"`
	AdminCountry string `gorm:"default:''" json:"admin_country"`
	
	// Tech Contact Info (defaults to Registrant if empty)
	TechName    string `gorm:"default:''" json:"tech_name"`
	TechOrg     string `gorm:"default:''" json:"tech_org"`
	TechEmail   string `gorm:"default:''" json:"tech_email"`
	TechPhone   string `gorm:"default:''" json:"tech_phone"`
	TechAddress string `gorm:"default:''" json:"tech_address"`
	TechCity    string `gorm:"default:''" json:"tech_city"`
	TechState   string `gorm:"default:''" json:"tech_state"`
	TechZip     string `gorm:"default:''" json:"tech_zip"`
	TechCountry string `gorm:"default:''" json:"tech_country"`
	
	// Status
	Status string `gorm:"default:'active'" json:"status"` // active, expired, suspended
	
	// Relations
	Records []Record `json:"records,omitempty"`
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

// RegistrarConfig stores global registrar settings
type RegistrarConfig struct {
	ID                uint   `gorm:"primaryKey" json:"id"`
	RegistrarName     string `gorm:"not null" json:"registrar_name"`
	RegistrarURL      string `json:"registrar_url"`
	RegistrarEmail    string `json:"registrar_email"`
	RegistrarPhone    string `json:"registrar_phone"`
	RegistrarIANAID   string `gorm:"default:'9999'" json:"registrar_iana_id"`
	AbuseContactEmail string `json:"abuse_contact_email"`
	AbuseContactPhone string `json:"abuse_contact_phone"`
	WhoisServer       string `json:"whois_server"`
	NameServer1       string `json:"nameserver1"`
	NameServer2       string `json:"nameserver2"`
	DefaultTTL        int    `gorm:"default:3600" json:"default_ttl"`
	DefaultExpiry     int    `gorm:"default:365" json:"default_expiry_days"` // Days until expiry
}
