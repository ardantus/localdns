package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/localdns/backend/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var SecretKey = []byte(os.Getenv("JWT_SECRET"))

func Register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 14)

		user := models.User{
			Username:     input.Username,
			PasswordHash: string(hashedPassword),
		}

		if result := db.Create(&user); result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username probably exists"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
	}
}

func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user models.User
		if result := db.Where("username = ?", input.Username).First(&user); result.Error != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// Generate JWT
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"iss":  "localdns",
			"sub":  user.ID,
			"role": user.Role,
			"exp":  time.Now().Add(time.Hour * 24).Unix(),
		})

		tokenString, err := token.SignedString(SecretKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": tokenString,
			"user":  gin.H{"id": user.ID, "username": user.Username, "role": user.Role},
		})
	}
}
