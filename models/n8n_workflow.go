package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type N8NWorkflow struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	UserID       uint           `gorm:"not null;index" json:"user_id"`
	WorkflowID   string         `gorm:"type:varchar(255)" json:"workflow_id"`
	Name         string         `gorm:"type:varchar(255);not null" json:"name"`
	Description  string         `gorm:"type:text" json:"description"`
	WorkflowJSON datatypes.JSON `gorm:"type:jsonb" json:"workflow_json"`
	ChatHistory  datatypes.JSON `gorm:"type:jsonb;not null" json:"chat_history"`
	Active       bool           `gorm:"default:false" json:"active"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	User         Users          `gorm:"foreignKey:UserID" json:"-"`
}

type WorkflowChatMessage struct {
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}
