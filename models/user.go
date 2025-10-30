package models

import (
	"time"

	"gorm.io/gorm"
)

type Users struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	FbID      string         `gorm:"uniqueIndex;not null" json:"fb_id"`
	Email     string         `gorm:"type:varchar(255)" json:"email"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
