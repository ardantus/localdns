package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/localdns/backend/models"
	"gorm.io/gorm"
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

		domain := models.Domain{
			Name:   input.Name,
			UserID: targetUserID,
		}

		if result := db.Create(&domain); result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Domain already exists or invalid"})
			return
		}

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
            // User sees only own domains
		    if result := db.Where("user_id = ?", userID).Find(&domains); result.Error != nil {
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
