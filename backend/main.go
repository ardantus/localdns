package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/localdns/backend/handlers"
    "github.com/localdns/backend/models"
    "golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// Wait loop or just fail to let Docker restart? Ideally wait, but for MVP we fail and restart.
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

    // Hack to fix "constraint does not exist" error on restart
    // Check if table exists, if so try to drop the conflicting index/constraint
    if db.Migrator().HasTable(&models.User{}) {
        // Ignore error if it doesn't exist
        db.Exec("ALTER TABLE users DROP CONSTRAINT IF EXISTS uni_users_username")
        db.Exec("DROP INDEX IF EXISTS uni_users_username")
    }

    // Auto-migrate tables individually to prevent one failure from blocking others
    if err := db.AutoMigrate(&models.RegistrarConfig{}); err != nil {
         log.Printf("Failed to auto-migrate RegistrarConfig: %v", err)
    }
    if err := db.AutoMigrate(&models.Domain{}, &models.Record{}); err != nil {
         log.Printf("Failed to auto-migrate Domain/Record: %v", err)
    }
    
    // User migration often fails on constraints, so we try soft migration then manual column headers
    if err := db.AutoMigrate(&models.User{}); err != nil {
        log.Printf("Failed to auto-migrate User: %v", err)
    }
    // Manual fallback for User columns if migration failed
    if db.Migrator().HasTable(&models.User{}) {
        m := db.Migrator()
        if !m.HasColumn(&models.User{}, "ContactName") { m.AddColumn(&models.User{}, "ContactName") }
        if !m.HasColumn(&models.User{}, "ContactOrg") { m.AddColumn(&models.User{}, "ContactOrg") }
        if !m.HasColumn(&models.User{}, "ContactEmail") { m.AddColumn(&models.User{}, "ContactEmail") }
        if !m.HasColumn(&models.User{}, "ContactPhone") { m.AddColumn(&models.User{}, "ContactPhone") }
        if !m.HasColumn(&models.User{}, "ContactAddress") { m.AddColumn(&models.User{}, "ContactAddress") }
        if !m.HasColumn(&models.User{}, "ContactCity") { m.AddColumn(&models.User{}, "ContactCity") }
        if !m.HasColumn(&models.User{}, "ContactState") { m.AddColumn(&models.User{}, "ContactState") }
        if !m.HasColumn(&models.User{}, "ContactZip") { m.AddColumn(&models.User{}, "ContactZip") }
        if !m.HasColumn(&models.User{}, "ContactCountry") { m.AddColumn(&models.User{}, "ContactCountry") }
    }

    // Manual fallback for Domain columns if migration failed
    if db.Migrator().HasTable(&models.Domain{}) {
        m := db.Migrator()
        // Core fields
        if !m.HasColumn(&models.Domain{}, "UpdatedAt") { m.AddColumn(&models.Domain{}, "UpdatedAt") }
        if !m.HasColumn(&models.Domain{}, "ExpiresAt") { m.AddColumn(&models.Domain{}, "ExpiresAt") }
        // Registrant fields
        if !m.HasColumn(&models.Domain{}, "RegistrantName") { m.AddColumn(&models.Domain{}, "RegistrantName") }
        if !m.HasColumn(&models.Domain{}, "RegistrantOrg") { m.AddColumn(&models.Domain{}, "RegistrantOrg") }
        if !m.HasColumn(&models.Domain{}, "RegistrantEmail") { m.AddColumn(&models.Domain{}, "RegistrantEmail") }
        if !m.HasColumn(&models.Domain{}, "RegistrantPhone") { m.AddColumn(&models.Domain{}, "RegistrantPhone") }
        if !m.HasColumn(&models.Domain{}, "RegistrantAddress") { m.AddColumn(&models.Domain{}, "RegistrantAddress") }
        if !m.HasColumn(&models.Domain{}, "RegistrantCity") { m.AddColumn(&models.Domain{}, "RegistrantCity") }
        if !m.HasColumn(&models.Domain{}, "RegistrantState") { m.AddColumn(&models.Domain{}, "RegistrantState") }
        if !m.HasColumn(&models.Domain{}, "RegistrantZip") { m.AddColumn(&models.Domain{}, "RegistrantZip") }
        if !m.HasColumn(&models.Domain{}, "RegistrantCountry") { m.AddColumn(&models.Domain{}, "RegistrantCountry") }
        // Admin fields
        if !m.HasColumn(&models.Domain{}, "AdminName") { m.AddColumn(&models.Domain{}, "AdminName") }
        if !m.HasColumn(&models.Domain{}, "AdminOrg") { m.AddColumn(&models.Domain{}, "AdminOrg") }
        if !m.HasColumn(&models.Domain{}, "AdminEmail") { m.AddColumn(&models.Domain{}, "AdminEmail") }
        if !m.HasColumn(&models.Domain{}, "AdminPhone") { m.AddColumn(&models.Domain{}, "AdminPhone") }
        if !m.HasColumn(&models.Domain{}, "AdminAddress") { m.AddColumn(&models.Domain{}, "AdminAddress") }
        if !m.HasColumn(&models.Domain{}, "AdminCity") { m.AddColumn(&models.Domain{}, "AdminCity") }
        if !m.HasColumn(&models.Domain{}, "AdminState") { m.AddColumn(&models.Domain{}, "AdminState") }
        if !m.HasColumn(&models.Domain{}, "AdminZip") { m.AddColumn(&models.Domain{}, "AdminZip") }
        if !m.HasColumn(&models.Domain{}, "AdminCountry") { m.AddColumn(&models.Domain{}, "AdminCountry") }
        // Tech fields
        if !m.HasColumn(&models.Domain{}, "TechName") { m.AddColumn(&models.Domain{}, "TechName") }
        if !m.HasColumn(&models.Domain{}, "TechOrg") { m.AddColumn(&models.Domain{}, "TechOrg") }
        if !m.HasColumn(&models.Domain{}, "TechEmail") { m.AddColumn(&models.Domain{}, "TechEmail") }
        if !m.HasColumn(&models.Domain{}, "TechPhone") { m.AddColumn(&models.Domain{}, "TechPhone") }
        if !m.HasColumn(&models.Domain{}, "TechAddress") { m.AddColumn(&models.Domain{}, "TechAddress") }
        if !m.HasColumn(&models.Domain{}, "TechCity") { m.AddColumn(&models.Domain{}, "TechCity") }
        if !m.HasColumn(&models.Domain{}, "TechState") { m.AddColumn(&models.Domain{}, "TechState") }
        if !m.HasColumn(&models.Domain{}, "TechZip") { m.AddColumn(&models.Domain{}, "TechZip") }
        if !m.HasColumn(&models.Domain{}, "TechCountry") { m.AddColumn(&models.Domain{}, "TechCountry") }
        // Status
        if !m.HasColumn(&models.Domain{}, "Status") { m.AddColumn(&models.Domain{}, "Status") }
    }

    // Seed Admin User
    var existingAdmin models.User
    if err := db.Where("username = ?", "admin").First(&existingAdmin).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            // Create default admin: admin / admin123
            hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), 14)
            admin := models.User{
                Username:     "admin",
                PasswordHash: string(hash),
                Role:         "admin",
            }
            if result := db.Create(&admin); result.Error != nil {
                log.Printf("Failed to seed admin user: %v", result.Error)
            } else {
                log.Println("Seeded default admin user: admin / admin123")
            }
        }
    }

    // Seed default Registrar Config
    var existingConfig models.RegistrarConfig
    if err := db.First(&existingConfig).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            config := models.RegistrarConfig{
                RegistrarName:     "LocalDNS Registrar",
                RegistrarURL:      "http://localhost:3000",
                RegistrarEmail:    "admin@localdns.local",
                RegistrarPhone:    "+1-555-0100",
                RegistrarIANAID:   "9999",
                AbuseContactEmail: "abuse@localdns.local",
                AbuseContactPhone: "+1-555-0199",
                WhoisServer:       "whois.localdns.local",
                NameServer1:       "ns1.localdns.local",
                NameServer2:       "ns2.localdns.local",
                DefaultTTL:        3600,
                DefaultExpiry:     365,
            }
            db.Create(&config)
            log.Println("Seeded default registrar configuration")
        }
    }


	r := gin.Default()

	// Public
	r.POST("/api/register", handlers.Register(db))
	r.POST("/api/login", handlers.Login(db))
	
	// Public WHOIS endpoint (no auth required)
	r.GET("/whois/:domain", handlers.WhoisRaw(db))
	r.GET("/api/whois", handlers.WhoisQuery(db))

	// Protected (TODO: Add Auth Middleware)
	api := r.Group("/api")
    api.Use(handlers.AuthMiddleware())
	{
		// Domains
		api.GET("/domains", handlers.ListDomains(db))
		api.POST("/domains", handlers.CreateDomain(db))
		api.GET("/domains/:id", handlers.GetDomain(db))
		api.DELETE("/domains/:id", handlers.DeleteDomain(db))
		api.PUT("/domains/:id/registrant", handlers.UpdateDomainRegistrant(db))
		
		// Records
		api.GET("/domains/:id/records", handlers.ListRecords(db))
		api.POST("/domains/:id/records", handlers.AddRecord(db))
		api.PUT("/records/:recordId", handlers.UpdateRecord(db))
		api.DELETE("/records/:recordId", handlers.DeleteRecord(db))
		
		// Users (admin only)
		api.GET("/users", handlers.ListUsers(db))
		api.POST("/users", handlers.CreateUser(db))
		api.PUT("/users/:id", handlers.UpdateUser(db))
		api.DELETE("/users/:id", handlers.DeleteUser(db))
		
		// Registrar Config (admin only for update)
		api.GET("/config", handlers.GetRegistrarConfig(db))
		api.PUT("/config", handlers.UpdateRegistrarConfig(db))
	}


	r.Run(":8080")
}
