package handler

import (
	"main/repository"
	"main/services"
)

type Handler struct {
	repo    *repository.Repository
	service *services.Service
}

func NewHandler(service *services.Service, repo *repository.Repository) *Handler {
	return &Handler{service: service, repo: repo}
}
