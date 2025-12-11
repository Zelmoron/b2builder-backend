package database

import (
	"log"
	"main/models"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDatabase(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
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
	log.Println("Users migration completed")

	if err := db.AutoMigrate(&models.Bot{}); err != nil {
		log.Fatal("Failed to migrate Bot:", err)
	}
	log.Println("Bot migration completed")

	if err := db.AutoMigrate(&models.ChatSession{}); err != nil {
		log.Fatal("Failed to migrate ChatSession:", err)
	}
	log.Println("ChatSession migration completed")

	log.Println("Database migration completed")
}
