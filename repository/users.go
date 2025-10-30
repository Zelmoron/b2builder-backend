package repository

import (
	"main/models"

	"gorm.io/gorm"
)

func (r *Repository) GetAllUsers() ([]models.Users, error) {
	var users []models.Users
	err := r.db.Find(&users).Error
	return users, err
}

func (r *Repository) CreateOrGetUser(fbID, email string) (*models.Users, error) {
	var user models.Users
	err := r.db.Where("fb_id = ?", fbID).First(&user).Error

	if err == gorm.ErrRecordNotFound {
		user = models.Users{
			FbID:  fbID,
			Email: email,
		}
		err = r.db.Create(&user).Error
		return &user, err
	}

	return &user, err
}

func (r *Repository) GetUserByFbID(fbID string) (*models.Users, error) {
	var user models.Users
	err := r.db.Where("fb_id = ?", fbID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
