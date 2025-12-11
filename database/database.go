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
	// Migrate Users - skip if already exists and has issues
	if err := db.AutoMigrate(&models.Users{}); err != nil {
		log.Println("Warning: Users migration had issues:", err)
	} else {
		log.Println("Users migration completed")
	}

	// Migrate Bot - skip if already exists and has issues
	if err := db.AutoMigrate(&models.Bot{}); err != nil {
		log.Println("Warning: Bot migration had issues:", err)
	} else {
		log.Println("Bot migration completed")
	}

	// Migrate ChatSession - skip if already exists and has issues
	if err := db.AutoMigrate(&models.ChatSession{}); err != nil {
		log.Println("Warning: ChatSession migration had issues:", err)
	} else {
		log.Println("ChatSession migration completed")
	}

	// Migrate N8NWorkflow
	if err := db.AutoMigrate(&models.N8NWorkflow{}); err != nil {
		log.Println("Warning: N8NWorkflow migration had issues:", err)
	} else {
		log.Println("N8NWorkflow migration completed")
	}

	log.Println("Database migration completed")
}
