package database

import (
	"log"
	"main/models"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to connect database")
	}
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetMaxOpenConns(5)
	sqlDB.SetConnMaxLifetime(time.Minute)
	return db
}

func Migrate(db *gorm.DB) {
	if err := db.AutoMigrate(&models.Users{}); err != nil {
		log.Fatal("Failed to migrate Users:", err)
	}

	if err := db.AutoMigrate(&models.Bot{}); err != nil {
		log.Fatal("Failed to migrate Bot:", err)
	}

	if err := db.AutoMigrate(&models.ChatSession{}); err != nil {
		log.Fatal("Failed to migrate ChatSession:", err)
	}

	log.Println("Database migration completed")
}
