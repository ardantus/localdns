package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Models for WHOIS server
type Domain struct {
	ID                uint `gorm:"primaryKey"`
	Name              string
	UserID            uint
	CreatedAt         time.Time
	UpdatedAt         time.Time
	ExpiresAt         time.Time
	RegistrantName    string
	RegistrantOrg     string
	RegistrantEmail   string
	RegistrantPhone   string
	RegistrantAddress string
	RegistrantCity    string
	RegistrantState   string
	RegistrantZip     string
	RegistrantCountry string
	AdminName         string
	AdminOrg          string
	AdminEmail        string
	AdminPhone        string
	AdminAddress      string
	AdminCity         string
	AdminState        string
	AdminZip          string
	AdminCountry      string
	TechName          string
	TechOrg           string
	TechEmail         string
	TechPhone         string
	TechAddress       string
	TechCity          string
	TechState         string
	TechZip           string
	TechCountry       string
	Status            string
}

type RegistrarConfig struct {
	ID                uint   `gorm:"primaryKey"`
	RegistrarName     string
	RegistrarURL      string
	RegistrarEmail    string
	RegistrarPhone    string
	RegistrarIANAID   string
	AbuseContactEmail string
	AbuseContactPhone string
	WhoisServer       string
	NameServer1       string
	NameServer2       string
}

var db *gorm.DB

func main() {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		getEnv("DB_HOST", "postgres"),
		getEnv("DB_USER", "user"),
		getEnv("DB_PASSWORD", "password"),
		getEnv("DB_NAME", "localdns"),
		getEnv("DB_PORT", "5432"),
	)

	for i := 0; i < 30; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Waiting for database... (%d/30)", i+1)
		time.Sleep(time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("WHOIS server starting on port 43...")
	listener, err := net.Listen("tcp", ":43")
	if err != nil {
		log.Fatalf("Failed to start WHOIS server: %v", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	query, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	domainName := strings.TrimSpace(strings.ToLower(query))
	log.Printf("WHOIS query for: %s", domainName)
	response := lookupDomain(domainName)
	conn.Write([]byte(response))
}

func lookupDomain(domainName string) string {
	var domain Domain
	if result := db.Where("name = ?", domainName).First(&domain); result.Error != nil {
		return formatNotFound(domainName)
	}
	var config RegistrarConfig
	db.First(&config)
	return formatWhoisResponse(domain, config)
}

func formatWhoisResponse(domain Domain, config RegistrarConfig) string {
	expiryDate := domain.ExpiresAt
	if expiryDate.IsZero() {
		expiryDate = domain.UpdatedAt.AddDate(1, 0, 0)
		if domain.UpdatedAt.IsZero() {
			expiryDate = domain.CreatedAt.AddDate(1, 0, 0)
		}
	}

	adminName := valueOrFallback(domain.AdminName, domain.RegistrantName)
	adminOrg := valueOrFallback(domain.AdminOrg, domain.RegistrantOrg)
	adminEmail := valueOrFallback(domain.AdminEmail, domain.RegistrantEmail)
	adminPhone := valueOrFallback(domain.AdminPhone, domain.RegistrantPhone)
	adminAddress := valueOrFallback(domain.AdminAddress, domain.RegistrantAddress)
	adminCity := valueOrFallback(domain.AdminCity, domain.RegistrantCity)
	adminState := valueOrFallback(domain.AdminState, domain.RegistrantState)
	adminZip := valueOrFallback(domain.AdminZip, domain.RegistrantZip)
	adminCountry := valueOrFallback(domain.AdminCountry, domain.RegistrantCountry)

	techName := valueOrFallback(domain.TechName, domain.RegistrantName)
	techOrg := valueOrFallback(domain.TechOrg, domain.RegistrantOrg)
	techEmail := valueOrFallback(domain.TechEmail, domain.RegistrantEmail)
	techPhone := valueOrFallback(domain.TechPhone, domain.RegistrantPhone)
	techAddress := valueOrFallback(domain.TechAddress, domain.RegistrantAddress)
	techCity := valueOrFallback(domain.TechCity, domain.RegistrantCity)
	techState := valueOrFallback(domain.TechState, domain.RegistrantState)
	techZip := valueOrFallback(domain.TechZip, domain.RegistrantZip)
	techCountry := valueOrFallback(domain.TechCountry, domain.RegistrantCountry)

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
		domain.ID,
		valueOrDefault(domain.RegistrantName, "REDACTED FOR PRIVACY"),
		valueOrDefault(domain.RegistrantOrg, "REDACTED FOR PRIVACY"),
		valueOrDefault(domain.RegistrantAddress, "REDACTED FOR PRIVACY"),
		valueOrDefault(domain.RegistrantCity, "REDACTED FOR PRIVACY"),
		valueOrDefault(domain.RegistrantState, "REDACTED FOR PRIVACY"),
		valueOrDefault(domain.RegistrantZip, "REDACTED FOR PRIVACY"),
		valueOrDefault(domain.RegistrantCountry, "REDACTED FOR PRIVACY"),
		valueOrDefault(domain.RegistrantPhone, "REDACTED FOR PRIVACY"),
		valueOrDefault(domain.RegistrantEmail, "REDACTED FOR PRIVACY"),
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
		config.NameServer1,
		config.NameServer2,
		time.Now().Format(time.RFC3339),
	)
}

func formatNotFound(domain string) string {
	return fmt.Sprintf(`No match for domain "%s".

>>> Last update of WHOIS database: %s <<<

NOTICE: This is a LOCAL DNS REGISTRAR for homelab/internal network use only.
`, domain, time.Now().Format(time.RFC3339))
}

func valueOrDefault(val, def string) string {
	if val == "" {
		return def
	}
	return val
}

func valueOrFallback(val, fallback string) string {
	if val == "" {
		return fallback
	}
	return val
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
