package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ChatSession struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	SessionID string         `gorm:"not null;index" json:"session_id"`
	BotID     string         `gorm:"not null;index" json:"bot_id"`
	Messages  datatypes.JSON `gorm:"type:jsonb;not null" json:"messages"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type ChatMessage struct {
	Role      string    `json:"role"` // "user" or "assistant"
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}
