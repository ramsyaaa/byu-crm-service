package service

import (
	"byu-crm-service/modules/territory/repository"
)

type territoryService struct {
	repo repository.TerritoryRepository
}

func NewTerritoryService(repo repository.TerritoryRepository) TerritoryService {
	return &territoryService{repo: repo}
}

func (s *territoryService) GetAllTerritories() (map[string]interface{}, int64, error) {
	return s.repo.GetAllTerritories()
}
