package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/localdns/backend/models"
	"gorm.io/gorm"
)

// WhoisQuery handles WHOIS lookups via HTTP API
func WhoisQuery(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		domainName := c.Query("domain")
		if domainName == "" {
			domainName = c.Param("domain")
		}
		if domainName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Domain name required"})
			return
		}

		domainName = strings.TrimSpace(strings.ToLower(domainName))

		var domain models.Domain
		if result := db.Where("name = ?", domainName).First(&domain); result.Error != nil {
			c.String(http.StatusNotFound, formatWhoisNotFound(domainName))
			return
		}

		// Get owner user for contact info
		var user models.User
		db.First(&user, domain.UserID)

		var config models.RegistrarConfig
		db.First(&config)

		c.String(http.StatusOK, formatWhoisResponse(domain, user, config))
	}
}

// WhoisRaw returns raw text WHOIS format
func WhoisRaw(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		domainName := c.Param("domain")
		domainName = strings.TrimSpace(strings.ToLower(domainName))

		var domain models.Domain
		if result := db.Where("name = ?", domainName).First(&domain); result.Error != nil {
			c.String(http.StatusNotFound, formatWhoisNotFound(domainName))
			return
		}

		// Get owner user for contact info
		var user models.User
		db.First(&user, domain.UserID)

		var config models.RegistrarConfig
		db.First(&config)

		c.Header("Content-Type", "text/plain; charset=utf-8")
		c.String(http.StatusOK, formatWhoisResponse(domain, user, config))
	}
}

func formatWhoisResponse(domain models.Domain, user models.User, config models.RegistrarConfig) string {
	// Calculate expiry: 1 year after last update
	expiryDate := domain.ExpiresAt
	if expiryDate.IsZero() {
		expiryDate = domain.UpdatedAt.AddDate(1, 0, 0)
		if domain.UpdatedAt.IsZero() {
			expiryDate = domain.CreatedAt.AddDate(1, 0, 0)
		}
	}

	// Use User contact data as fallback for Registrant
	registrantName := valueOrFallback(domain.RegistrantName, user.ContactName)
	registrantOrg := valueOrFallback(domain.RegistrantOrg, user.ContactOrg)
	registrantEmail := valueOrFallback(domain.RegistrantEmail, user.ContactEmail)
	registrantPhone := valueOrFallback(domain.RegistrantPhone, user.ContactPhone)
	registrantAddress := valueOrFallback(domain.RegistrantAddress, user.ContactAddress)
	registrantCity := valueOrFallback(domain.RegistrantCity, user.ContactCity)
	registrantState := valueOrFallback(domain.RegistrantState, user.ContactState)
	registrantZip := valueOrFallback(domain.RegistrantZip, user.ContactZip)
	registrantCountry := valueOrFallback(domain.RegistrantCountry, user.ContactCountry)

	// Admin falls back to Registrant, then User
	adminName := valueOrFallback(domain.AdminName, valueOrFallback(registrantName, user.ContactName))
	adminOrg := valueOrFallback(domain.AdminOrg, valueOrFallback(registrantOrg, user.ContactOrg))
	adminEmail := valueOrFallback(domain.AdminEmail, valueOrFallback(registrantEmail, user.ContactEmail))
	adminPhone := valueOrFallback(domain.AdminPhone, valueOrFallback(registrantPhone, user.ContactPhone))
	adminAddress := valueOrFallback(domain.AdminAddress, valueOrFallback(registrantAddress, user.ContactAddress))
	adminCity := valueOrFallback(domain.AdminCity, valueOrFallback(registrantCity, user.ContactCity))
	adminState := valueOrFallback(domain.AdminState, valueOrFallback(registrantState, user.ContactState))
	adminZip := valueOrFallback(domain.AdminZip, valueOrFallback(registrantZip, user.ContactZip))
	adminCountry := valueOrFallback(domain.AdminCountry, valueOrFallback(registrantCountry, user.ContactCountry))

	// Tech falls back to Registrant, then User
	techName := valueOrFallback(domain.TechName, valueOrFallback(registrantName, user.ContactName))
	techOrg := valueOrFallback(domain.TechOrg, valueOrFallback(registrantOrg, user.ContactOrg))
	techEmail := valueOrFallback(domain.TechEmail, valueOrFallback(registrantEmail, user.ContactEmail))
	techPhone := valueOrFallback(domain.TechPhone, valueOrFallback(registrantPhone, user.ContactPhone))
	techAddress := valueOrFallback(domain.TechAddress, valueOrFallback(registrantAddress, user.ContactAddress))
	techCity := valueOrFallback(domain.TechCity, valueOrFallback(registrantCity, user.ContactCity))
	techState := valueOrFallback(domain.TechState, valueOrFallback(registrantState, user.ContactState))
	techZip := valueOrFallback(domain.TechZip, valueOrFallback(registrantZip, user.ContactZip))
	techCountry := valueOrFallback(domain.TechCountry, valueOrFallback(registrantCountry, user.ContactCountry))


	return fmt.Sprintf(`Domain Name: %s
Registry Domain ID: DOM-%d-LOCALDNS
Registrar WHOIS Server: %s
Registrar URL: %s
Updated Date: %s
Creation Date: %s
Registry Expiry Date: %s
Registrar: %s
Registrar IANA ID: %s
Registrar Abuse Contact Email: %s
Registrar Abuse Contact Phone: %s
Domain Status: %s https://icann.org/epp#%s

Registry Registrant ID: C%d-LOCALDNS
Registrant Name: %s
Registrant Organization: %s
Registrant Street: %s
Registrant City: %s
Registrant State/Province: %s
Registrant Postal Code: %s
Registrant Country: %s
Registrant Phone: %s
Registrant Email: %s

Registry Admin ID: C%d-LOCALDNS
Admin Name: %s
Admin Organization: %s
Admin Street: %s
Admin City: %s
Admin State/Province: %s
Admin Postal Code: %s
Admin Country: %s
Admin Phone: %s
Admin Email: %s

Registry Tech ID: C%d-LOCALDNS
Tech Name: %s
Tech Organization: %s
Tech Street: %s
Tech City: %s
Tech State/Province: %s
Tech Postal Code: %s
Tech Country: %s
Tech Phone: %s
Tech Email: %s

Name Server: %s
Name Server: %s
DNSSEC: unsigned

>>> Last update of WHOIS database: %s <<<

TERMS OF USE: This WHOIS data is provided for informational purposes only.
This data conforms to RFC 3912 WHOIS protocol specification.

NOTICE: This is a LOCAL DNS REGISTRAR for homelab/internal network use only.
This WHOIS information follows IANA/ICANN formatting standards for educational
and testing purposes. This is NOT a real domain registration and has no legal
standing outside of your local network environment.

For more information on WHOIS status codes, please visit https://icann.org/epp

`,
		strings.ToUpper(domain.Name),
		domain.ID,
		config.WhoisServer,
		config.RegistrarURL,
		domain.UpdatedAt.Format(time.RFC3339),
		domain.CreatedAt.Format(time.RFC3339),
		expiryDate.Format(time.RFC3339),
		config.RegistrarName,
		valueOrDefault(config.RegistrarIANAID, "9999"),
		valueOrDefault(config.AbuseContactEmail, config.RegistrarEmail),
		valueOrDefault(config.AbuseContactPhone, config.RegistrarPhone),
		valueOrDefault(domain.Status, "active"),
		valueOrDefault(domain.Status, "active"),
		// Registrant
		domain.ID,
		valueOrDefault(registrantName, "REDACTED FOR PRIVACY"),
		valueOrDefault(registrantOrg, "REDACTED FOR PRIVACY"),
		valueOrDefault(registrantAddress, "REDACTED FOR PRIVACY"),
		valueOrDefault(registrantCity, "REDACTED FOR PRIVACY"),
		valueOrDefault(registrantState, "REDACTED FOR PRIVACY"),
		valueOrDefault(registrantZip, "REDACTED FOR PRIVACY"),
		valueOrDefault(registrantCountry, "REDACTED FOR PRIVACY"),
		valueOrDefault(registrantPhone, "REDACTED FOR PRIVACY"),
		valueOrDefault(registrantEmail, "REDACTED FOR PRIVACY"),
		// Admin
		domain.ID,
		valueOrDefault(adminName, "REDACTED FOR PRIVACY"),
		valueOrDefault(adminOrg, "REDACTED FOR PRIVACY"),
		valueOrDefault(adminAddress, "REDACTED FOR PRIVACY"),
		valueOrDefault(adminCity, "REDACTED FOR PRIVACY"),
		valueOrDefault(adminState, "REDACTED FOR PRIVACY"),
		valueOrDefault(adminZip, "REDACTED FOR PRIVACY"),
		valueOrDefault(adminCountry, "REDACTED FOR PRIVACY"),
		valueOrDefault(adminPhone, "REDACTED FOR PRIVACY"),
		valueOrDefault(adminEmail, "REDACTED FOR PRIVACY"),
		// Tech
		domain.ID,
		valueOrDefault(techName, "REDACTED FOR PRIVACY"),
		valueOrDefault(techOrg, "REDACTED FOR PRIVACY"),
		valueOrDefault(techAddress, "REDACTED FOR PRIVACY"),
		valueOrDefault(techCity, "REDACTED FOR PRIVACY"),
		valueOrDefault(techState, "REDACTED FOR PRIVACY"),
		valueOrDefault(techZip, "REDACTED FOR PRIVACY"),
		valueOrDefault(techCountry, "REDACTED FOR PRIVACY"),
		valueOrDefault(techPhone, "REDACTED FOR PRIVACY"),
		valueOrDefault(techEmail, "REDACTED FOR PRIVACY"),
		// Nameservers
		config.NameServer1,
		config.NameServer2,
		time.Now().Format(time.RFC3339),
	)
}

func formatWhoisNotFound(domain string) string {
	return fmt.Sprintf(`No match for domain "%s".

>>> Last update of WHOIS database: %s <<<

NOTICE: This is a LOCAL DNS REGISTRAR for homelab/internal network use only.
`, domain, time.Now().Format(time.RFC3339))
}

func valueOrDefault(val, def string) string {
	if strings.TrimSpace(val) == "" {
		return def
	}
	return val
}

func valueOrFallback(val, fallback string) string {
	if strings.TrimSpace(val) == "" {
		return fallback
	}
	return val
}

// GetRegistrarConfig returns the registrar configuration
func GetRegistrarConfig(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var config models.RegistrarConfig
		if result := db.First(&config); result.Error != nil {
			// If not found, return default empty config (or seed it)
			config.RegistrarName = "LocalDNS Registrar"
            config.RegistrarIANAID = "9999"
            config.DefaultTTL = 3600
            config.DefaultExpiry = 365
		}
		c.JSON(http.StatusOK, config)
	}
}

// UpdateRegistrarConfig updates the registrar configuration (admin only)
func UpdateRegistrarConfig(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.MustGet("role").(string)
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			return
		}

		var config models.RegistrarConfig
		if result := db.First(&config); result.Error != nil {
            // If not found, create new
            config = models.RegistrarConfig{}
        }

		var input struct {
			RegistrarName     string `json:"registrar_name"`
			RegistrarURL      string `json:"registrar_url"`
			RegistrarEmail    string `json:"registrar_email"`
			RegistrarPhone    string `json:"registrar_phone"`
			RegistrarIANAID   string `json:"registrar_iana_id"`
			AbuseContactEmail string `json:"abuse_contact_email"`
			AbuseContactPhone string `json:"abuse_contact_phone"`
			WhoisServer       string `json:"whois_server"`
			NameServer1       string `json:"nameserver1"`
			NameServer2       string `json:"nameserver2"`
			DefaultTTL        int    `json:"default_ttl"`
			DefaultExpiry     int    `json:"default_expiry_days"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update all fields directly
		config.RegistrarName = input.RegistrarName
		config.RegistrarURL = input.RegistrarURL
		config.RegistrarEmail = input.RegistrarEmail
		config.RegistrarPhone = input.RegistrarPhone
		config.RegistrarIANAID = input.RegistrarIANAID
		config.AbuseContactEmail = input.AbuseContactEmail
		config.AbuseContactPhone = input.AbuseContactPhone
		config.WhoisServer = input.WhoisServer
		config.NameServer1 = input.NameServer1
		config.NameServer2 = input.NameServer2
		if input.DefaultTTL > 0 {
			config.DefaultTTL = input.DefaultTTL
		}
		if input.DefaultExpiry > 0 {
			config.DefaultExpiry = input.DefaultExpiry
		}

		if err := db.Save(&config).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save config: " + err.Error()})
            return
        }
		c.JSON(http.StatusOK, config)
	}
}

// UpdateDomainRegistrant updates registrant info for a domain
func UpdateDomainRegistrant(db *gorm.DB) gin.HandlerFunc {
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

		var input struct {
			// Registrant
			RegistrantName    string `json:"registrant_name"`
			RegistrantOrg     string `json:"registrant_org"`
			RegistrantEmail   string `json:"registrant_email"`
			RegistrantPhone   string `json:"registrant_phone"`
			RegistrantAddress string `json:"registrant_address"`
			RegistrantCity    string `json:"registrant_city"`
			RegistrantState   string `json:"registrant_state"`
			RegistrantZip     string `json:"registrant_zip"`
			RegistrantCountry string `json:"registrant_country"`
			// Admin
			AdminName    string `json:"admin_name"`
			AdminOrg     string `json:"admin_org"`
			AdminEmail   string `json:"admin_email"`
			AdminPhone   string `json:"admin_phone"`
			AdminAddress string `json:"admin_address"`
			AdminCity    string `json:"admin_city"`
			AdminState   string `json:"admin_state"`
			AdminZip     string `json:"admin_zip"`
			AdminCountry string `json:"admin_country"`
			// Tech
			TechName    string `json:"tech_name"`
			TechOrg     string `json:"tech_org"`
			TechEmail   string `json:"tech_email"`
			TechPhone   string `json:"tech_phone"`
			TechAddress string `json:"tech_address"`
			TechCity    string `json:"tech_city"`
			TechState   string `json:"tech_state"`
			TechZip     string `json:"tech_zip"`
			TechCountry string `json:"tech_country"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update Registrant
		domain.RegistrantName = input.RegistrantName
		domain.RegistrantOrg = input.RegistrantOrg
		domain.RegistrantEmail = input.RegistrantEmail
		domain.RegistrantPhone = input.RegistrantPhone
		domain.RegistrantAddress = input.RegistrantAddress
		domain.RegistrantCity = input.RegistrantCity
		domain.RegistrantState = input.RegistrantState
		domain.RegistrantZip = input.RegistrantZip
		domain.RegistrantCountry = input.RegistrantCountry
		// Update Admin
		domain.AdminName = input.AdminName
		domain.AdminOrg = input.AdminOrg
		domain.AdminEmail = input.AdminEmail
		domain.AdminPhone = input.AdminPhone
		domain.AdminAddress = input.AdminAddress
		domain.AdminCity = input.AdminCity
		domain.AdminState = input.AdminState
		domain.AdminZip = input.AdminZip
		domain.AdminCountry = input.AdminCountry
		// Update Tech
		domain.TechName = input.TechName
		domain.TechOrg = input.TechOrg
		domain.TechEmail = input.TechEmail
		domain.TechPhone = input.TechPhone
		domain.TechAddress = input.TechAddress
		domain.TechCity = input.TechCity
		domain.TechState = input.TechState
		domain.TechZip = input.TechZip
		domain.TechCountry = input.TechCountry

		// Update expiry to +1 year from now on any update
		domain.ExpiresAt = time.Now().AddDate(1, 0, 0)

		db.Save(&domain)
		c.JSON(http.StatusOK, domain)
	}
}

// GetDomain returns a single domain with all details
func GetDomain(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("user_id").(uint)
		role := c.MustGet("role").(string)
		domainID := c.Param("id")

		var domain models.Domain
		if result := db.Preload("User").Preload("Records").First(&domain, domainID); result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Domain not found"})
			return
		}

		if role != "admin" && domain.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}

		c.JSON(http.StatusOK, domain)
	}
}
