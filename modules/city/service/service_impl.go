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

func (s *cityService) GetAllCities(limit int, paginate bool, page int, filters map[string]string, userRole string, territoryID int) ([]models.City, int64, error) {
	return s.repo.GetAllCities(limit, paginate, page, filters, userRole, territoryID)
}

func (s *cityService) GetCityByID(id uint) (*models.City, error) {
	return s.repo.FindByID(id)
}

func (s *cityService) GetCityByName(name string) (*models.City, error) {
	return s.repo.FindByName(name)
}
