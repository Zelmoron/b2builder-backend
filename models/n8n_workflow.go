package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type N8NWorkflow struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	UserID       uint           `gorm:"not null;index" json:"user_id"` // Internal user ID (references Users.ID)
	WorkflowID   string         `gorm:"type:varchar(255)" json:"workflow_id"` // N8N workflow ID (can be null if not yet created)
	Name         string         `gorm:"type:varchar(255);not null" json:"name"` // User-defined workflow name
	Description  string         `gorm:"type:text" json:"description"` // Workflow description
	WorkflowJSON datatypes.JSON `gorm:"type:jsonb" json:"workflow_json"` // Current N8N workflow JSON
	ChatHistory  datatypes.JSON `gorm:"type:jsonb;not null" json:"chat_history"` // Conversation history with AI
	Active       bool           `gorm:"default:false" json:"active"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// ChatMessage represents a single message in the conversation
type WorkflowChatMessage struct {
	Role      string    `json:"role"` // "user" or "assistant"
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}