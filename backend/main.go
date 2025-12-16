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

    // Auto-migrate
    db.AutoMigrate(&models.User{}, &models.Domain{}, &models.Record{})

    // Seed Admin User
    var adminCount int64
    db.Model(&models.User{}).Where("role = ?", "admin").Count(&adminCount)
    if adminCount == 0 {
        // Create default admin: admin / admin123
        // Hash for "admin123"
        hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), 14)
        admin := models.User{
            Username: "admin",
            PasswordHash: string(hash),
            Role: "admin",
        }
        db.Create(&admin)
        log.Println("Seeded default admin user: admin / admin123")
    }


	r := gin.Default()

	// Public
	r.POST("/api/register", handlers.Register(db))
	r.POST("/api/login", handlers.Login(db))

	// Protected (TODO: Add Auth Middleware)
	api := r.Group("/api")
    api.Use(handlers.AuthMiddleware())
	{
		api.GET("/domains", handlers.ListDomains(db))
		api.POST("/domains", handlers.CreateDomain(db))
		api.POST("/domains/:id/records", handlers.AddRecord(db))
	}


	r.Run(":8080")
}
