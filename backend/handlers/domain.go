package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/localdns/backend/models"
    "golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
    "strings"
)

func CreateDomain(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("user_id").(uint)
        role := c.MustGet("role").(string)

		var input struct {
			Name   string `json:"name" binding:"required"`
            UserID uint   `json:"user_id"` // Admin can specify user_id
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

        targetUserID := userID
        if role == "admin" && input.UserID != 0 {
            targetUserID = input.UserID
        }

		// Load user data to copy contact information
		var owner models.User
		if result := db.First(&owner, targetUserID); result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		// Get registrar config for default expiry
		var config models.RegistrarConfig
		defaultExpiryDays := 365 // Default fallback
		if err := db.First(&config).Error; err == nil && config.DefaultExpiry > 0 {
			defaultExpiryDays = config.DefaultExpiry
		}

		// Create domain with contact info from owner
		domain := models.Domain{
			Name:   input.Name,
			UserID: targetUserID,
			Status: "active",
			// Copy contact info from owner to registrant fields
			RegistrantName:    owner.ContactName,
			RegistrantOrg:     owner.ContactOrg,
			RegistrantEmail:   owner.ContactEmail,
			RegistrantPhone:   owner.ContactPhone,
			RegistrantAddress: owner.ContactAddress,
			RegistrantCity:    owner.ContactCity,
			RegistrantState:   owner.ContactState,
			RegistrantZip:     owner.ContactZip,
			RegistrantCountry: owner.ContactCountry,
			// Set expiry date
			ExpiresAt: time.Now().AddDate(0, 0, defaultExpiryDays),
		}

		if result := db.Create(&domain); result.Error != nil {
            // Check for unique constraint violation explicitly
            if strings.Contains(result.Error.Error(), "duplicate key value") || strings.Contains(result.Error.Error(), "UNIQUE constraint") {
			    c.JSON(http.StatusBadRequest, gin.H{"error": "Domain already exists"})
            } else {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create domain: " + result.Error.Error()})
            }
			return
		}

		// Reload domain with user association for response
		db.Preload("User").First(&domain, domain.ID)

		c.JSON(http.StatusCreated, domain)
	}
}

func ListDomains(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
        userID := c.MustGet("user_id").(uint)
        role := c.MustGet("role").(string)

		var domains []models.Domain
        
        if role == "admin" {
            // Admin sees all domains with owner info
            if result := db.Preload("User").Find(&domains); result.Error != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
                return
            }
        } else {
            // User sees only own domains, but still load user info for display
		    if result := db.Preload("User").Where("user_id = ?", userID).Find(&domains); result.Error != nil {
			    c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			    return
		    }
        }
        
		c.JSON(http.StatusOK, domains)
	}
}

func AddRecord(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
        userID := c.MustGet("user_id").(uint)
        role := c.MustGet("role").(string)

		domainID := c.Param("id")
		
		var input models.Record
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		// Ensure domain exists and owned by user (or is admin)
		var domain models.Domain
		if result := db.First(&domain, domainID); result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Domain not found"})
			return
		}

        if role != "admin" && domain.UserID != userID {
            c.JSON(http.StatusForbidden, gin.H{"error": "You do not own this domain"})
            return
        }

		input.DomainID = domain.ID
		// Force default if 0
		if input.TTL == 0 {
			input.TTL = 360
		}
		
		if result := db.Create(&input); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		c.JSON(http.StatusCreated, input)
	}
}

// ListRecords returns all records for a domain
func ListRecords(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("user_id").(uint)
		role := c.MustGet("role").(string)
		domainID := c.Param("id")

		var domain models.Domain
		if result := db.First(&domain, domainID); result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Domain not found"})
			return
		}

		if role != "admin" && domain.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}

		var records []models.Record
		db.Where("domain_id = ?", domain.ID).Find(&records)
		c.JSON(http.StatusOK, records)
	}
}

// DeleteRecord removes a DNS record
func DeleteRecord(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("user_id").(uint)
		role := c.MustGet("role").(string)
		recordID := c.Param("recordId")

		var record models.Record
		if result := db.First(&record, recordID); result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
			return
		}

		// Check ownership via domain
		var domain models.Domain
		db.First(&domain, record.DomainID)
		if role != "admin" && domain.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}

		db.Delete(&record)
		c.JSON(http.StatusOK, gin.H{"message": "Record deleted"})
	}
}

// DeleteDomain removes a domain and all its records
func DeleteDomain(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("user_id").(uint)
		role := c.MustGet("role").(string)
		domainID := c.Param("id")

		var domain models.Domain
		if result := db.First(&domain, domainID); result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Domain not found"})
			return
		}

		if role != "admin" && domain.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}

		// Delete all records first
		db.Where("domain_id = ?", domain.ID).Delete(&models.Record{})
		db.Delete(&domain)
		c.JSON(http.StatusOK, gin.H{"message": "Domain deleted"})
	}
}

// ListUsers returns all users (admin only)
func ListUsers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.MustGet("role").(string)
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			return
		}

		var users []models.User
		db.Find(&users)
		c.JSON(http.StatusOK, users)
	}
}

// DeleteUser removes a user (admin only)
func DeleteUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.MustGet("role").(string)
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			return
		}

		userID := c.Param("id")
		var user models.User
		if result := db.First(&user, userID); result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		// Don't allow deleting the last admin
		if user.Role == "admin" {
			var adminCount int64
			db.Model(&models.User{}).Where("role = ?", "admin").Count(&adminCount)
			if adminCount <= 1 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete the last admin"})
				return
			}
		}

		db.Delete(&user)
		c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
	}
}

// UpdateRecord modifies an existing DNS record
func UpdateRecord(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("user_id").(uint)
		role := c.MustGet("role").(string)
		recordID := c.Param("recordId")

		var record models.Record
		if result := db.First(&record, recordID); result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
			return
		}

		// Check ownership via domain
		var domain models.Domain
		db.First(&domain, record.DomainID)
		if role != "admin" && domain.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}

		var input struct {
			Name    string `json:"name"`
			Type    string `json:"type"`
			Content string `json:"content"`
			TTL     int    `json:"ttl"`
			Prio    int    `json:"prio"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update fields
		if input.Name != "" {
			record.Name = input.Name
		}
		if input.Type != "" {
			record.Type = input.Type
		}
		if input.Content != "" {
			record.Content = input.Content
		}
		if input.TTL > 0 {
			record.TTL = input.TTL
		}
		record.Prio = input.Prio

		db.Save(&record)
		c.JSON(http.StatusOK, record)
	}
}

// UpdateUser modifies an existing user (admin only)
func UpdateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		adminRole := c.MustGet("role").(string)
		if adminRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			return
		}

		targetUserID := c.Param("id")
		var user models.User
		if result := db.First(&user, targetUserID); result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		var input struct {
			Username       string `json:"username"`
			Role           string `json:"role"`
			ContactName    string `json:"contact_name"`
			ContactOrg     string `json:"contact_org"`
			ContactEmail   string `json:"contact_email"`
			ContactPhone   string `json:"contact_phone"`
			ContactAddress string `json:"contact_address"`
			ContactCity    string `json:"contact_city"`
			ContactState   string `json:"contact_state"`
			ContactZip     string `json:"contact_zip"`
			ContactCountry string `json:"contact_country"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Prevent demoting the last admin
		if user.Role == "admin" && input.Role == "user" {
			var adminCount int64
			db.Model(&models.User{}).Where("role = ?", "admin").Count(&adminCount)
			if adminCount <= 1 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot demote the last admin"})
				return
			}
		}

		if input.Username != "" {
			user.Username = input.Username
		}
		if input.Role != "" {
			user.Role = input.Role
		}
		// Update contact fields
		user.ContactName = input.ContactName
		user.ContactOrg = input.ContactOrg
		user.ContactEmail = input.ContactEmail
		user.ContactPhone = input.ContactPhone
		user.ContactAddress = input.ContactAddress
		user.ContactCity = input.ContactCity
		user.ContactState = input.ContactState
		user.ContactZip = input.ContactZip
		user.ContactCountry = input.ContactCountry

		if err := db.Save(&user).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user: " + err.Error()})
            return
        }
		c.JSON(http.StatusOK, user)
	}
}

// CreateUser creates a new user (admin only)
func CreateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleCtx := c.MustGet("role").(string)
		if roleCtx != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			return
		}

		var input struct {
			Username     string `json:"username" binding:"required"`
			Password     string `json:"password" binding:"required"`
			Role         string `json:"role"`
			ContactName    string `json:"contact_name"`
			ContactOrg     string `json:"contact_org"`
			ContactEmail   string `json:"contact_email"`
			ContactPhone   string `json:"contact_phone"`
			ContactAddress string `json:"contact_address"`
			ContactCity    string `json:"contact_city"`
			ContactState   string `json:"contact_state"`
			ContactZip     string `json:"contact_zip"`
			ContactCountry string `json:"contact_country"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

        // Check if user exists
        var count int64
        db.Model(&models.User{}).Where("username = ?", input.Username).Count(&count)
        if count > 0 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
            return
        }

        hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 14)
        if err != nil {
             c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
             return
        }

        user := models.User{
            Username:     input.Username,
            PasswordHash: string(hash),
            Role:         input.Role,
            ContactName:    input.ContactName,
            ContactOrg:     input.ContactOrg,
            ContactEmail:   input.ContactEmail,
            ContactPhone:   input.ContactPhone,
            ContactAddress: input.ContactAddress,
            ContactCity:    input.ContactCity,
            ContactState:   input.ContactState,
            ContactZip:     input.ContactZip,
            ContactCountry: input.ContactCountry,
        }
        if user.Role == "" {
             user.Role = "user"
        }

		if err := db.Create(&user).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
            return
        }
		c.JSON(http.StatusCreated, user)
	}
}

