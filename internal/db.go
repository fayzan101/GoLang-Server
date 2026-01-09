package internal

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func InitDB(connStr string) {
	var err error
	DB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Connected to PostgreSQL database with GORM.")

	// Auto-migrate User model
	type User struct {
		ID         uint   `gorm:"primaryKey"`
		FullName   string
		DOB        string
		University string
		Semester   string
		Program    string
		RollNo     string
		Email      string
		Password   string
		Type       string
	}
	if err := DB.AutoMigrate(&User{}); err != nil {
		log.Fatalf("Auto-migration failed: %v", err)
	}
}
