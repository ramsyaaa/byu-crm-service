package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/city/repository"
)

type cityService struct {
	repo repository.CityRepository
}

func NewCityService(repo repository.CityRepository) CityService {
	return &cityService{repo: repo}
}

func (s *cityService) GetCityByID(id uint) (*models.City, error) {
	return s.repo.FindByID(id)
}

func (s *cityService) GetCityByName(name string) (*models.City, error) {
	return s.repo.FindByName(name)
}
