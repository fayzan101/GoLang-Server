package internal

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(connStr string) {
	var err error
	DB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Connected to PostgreSQL database with GORM.")

	// Auto-migrate all IMS models
	if err := DB.AutoMigrate(
		&Product{},
		&Warehouse{},
		&Inventory{},
		&StockMovement{},
		&Supplier{},
		&PurchaseOrder{},
		&POItem{},
		&Order{},
		&OrderItem{},
		&AuditLog{},
	); err != nil {
		log.Fatalf("Auto-migration failed: %v", err)
	}
	log.Println("Database migration completed successfully.")
}

// LogAudit creates an audit log entry
func LogAudit(action, entity string, entityID uint, userID, details string) {
	log := AuditLog{
		Action:    action,
		Entity:    entity,
		EntityID:  entityID,
		UserID:    userID,
		Details:   details,
		IPAddress: "127.0.0.1",
		CreatedAt: time.Now(),
	}
	DB.Create(&log)
}
