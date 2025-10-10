package services

import (
	"main/models"
)

func (s *Service) GetAllUsers() ([]models.Users, error) {
	return s.repo.GetAllUsers()
}
