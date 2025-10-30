package services

import (
	"main/models"
)

func (s *Service) GetAllUsers() ([]models.Users, error) {
	return s.repo.GetAllUsers()
}

func (s *Service) RegisterUser(fbID, email string) (*models.Users, error) {
	return s.repo.CreateOrGetUser(fbID, email)
}
