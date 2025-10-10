package repository

import (
	"main/models"
)

func (r *Repository) GetAllUsers() ([]models.Users, error) {
	var users []models.Users
	err := r.db.Find(&users).Error
	return users, err
}
