package database

import (
	"log"
	"nnnn/main/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitializeDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	MigrateDB(db)

	return db
}

func MigrateDB(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{},
		&models.Notification{},
		// Add other models here
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}
