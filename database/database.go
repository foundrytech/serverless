package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID                        string     `gorm:"type:uuid;default:gen_random_uuid();primaryKey;->" json:"id"`
	FirstName                 string     `gorm:"not null" json:"first_name"`
	LastName                  string     `gorm:"not null" json:"last_name"`
	Password                  string     `gorm:"not null" json:"password"`
	Username                  string     `gorm:"unique;not null" json:"username"`
	Verified                  bool       `gorm:"default:false" json:"verified"`
	VerificationToken         *string    `gorm:"default:null" json:"verification_token"`
	VerificationTokenCreated  *time.Time `gorm:"default:null" json:"verification_token_created"`
	AccountCreated            time.Time  `gorm:"default:current_timestamp" json:"account_created"`
	AccountUpdated            time.Time  `gorm:"default:current_timestamp" json:"account_updated"`
}

var DB *gorm.DB

func Connect() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
	}
	log.Println("Connected to database")
}


func SaveTokenInfo(username string, token string, tokenCreated time.Time) error {
	user := User{}
	result := DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		log.Printf("Error querying user from database: %v", result.Error)
		return result.Error
	}

	// Update the VerificationToken and VerificationTokenCreated fields
	user.VerificationToken = &token
	user.VerificationTokenCreated = &tokenCreated

	result = DB.Save(&user)
	if result.Error != nil {
		log.Printf("Error saving token info to database: %v", result.Error)
		return result.Error
	}

	log.Printf("Token info saved to database: %v", user)
	return nil
}