package models

import (
	"time"

	"gorm.io/gorm"
)

type Bot struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	BotID        string         `gorm:"uniqueIndex;not null" json:"bot_id"`
	UserID       string         `gorm:"not null;index" json:"user_id"` // Firebase UID
	Name         string         `gorm:"type:varchar(255);not null" json:"name"`
	Type         string         `gorm:"type:varchar(50);not null" json:"type"`
	Status       string         `gorm:"type:varchar(20);default:'active'" json:"status"`
	SystemPrompt string         `gorm:"type:text" json:"system_prompt"` // Stores productDescription and FAQ as AI system prompt
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// FAQItem represents a single FAQ entry
type FAQItem struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}
