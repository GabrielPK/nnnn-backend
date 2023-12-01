package database

import (
    "gorm.io/gorm"
    "gorm.io/driver/sqlite"
    "log"
    "nnnn/main/models" 
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
        // Add other models here
    )
    if err != nil {
        log.Fatalf("Failed to migrate database: %v", err)
    }
}