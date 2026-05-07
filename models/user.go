package models

import "gorm.io/gorm"

type Users struct {
	gorm.Model
	FbID         string        `gorm:"uniqueIndex;not null" json:"fb_id"`
	Email        string        `gorm:"type:varchar(255)" json:"email"`
}