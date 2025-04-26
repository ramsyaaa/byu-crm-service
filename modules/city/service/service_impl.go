package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/city/repository"
	"byu-crm-service/modules/city/response"
	"fmt"
)

type cityService struct {
	repo repository.CityRepository
}

func NewCityService(repo repository.CityRepository) CityService {
	return &cityService{repo: repo}
}

func (s *cityService) GetAllCities(limit int, paginate bool, page int, filters map[string]string, userRole string, territoryID int) ([]response.CityResponse, int64, error) {
	return s.repo.GetAllCities(limit, paginate, page, filters, userRole, territoryID)
}

func (s *cityService) GetCityByID(id int) (*response.CityResponse, error) {
	return s.repo.GetCityByID(id)
}

func (s *cityService) GetCityByName(name string) (*response.CityResponse, error) {
	return s.repo.GetCityByName(name)
}

func (s *cityService) CreateCity(name *string, cluster_id int) (*response.CityResponse, error) {
	clusterIDStr := fmt.Sprintf("%d", cluster_id)
	city := &models.City{Name: *name, ClusterID: &clusterIDStr}
	return s.repo.CreateCity(city)
}

func (s *cityService) UpdateCity(name *string, cluster_id int, id int) (*response.CityResponse, error) {
	clusterIDStr := fmt.Sprintf("%d", cluster_id)
	city := &models.City{Name: *name, ClusterID: &clusterIDStr}
	return s.repo.UpdateCity(city, id)
}
