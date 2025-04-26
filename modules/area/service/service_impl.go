package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/area/repository"
	"byu-crm-service/modules/area/response"
)

type areaService struct {
	repo repository.AreaRepository
}

func NewAreaService(repo repository.AreaRepository) AreaService {
	return &areaService{repo: repo}
}

func (s *areaService) GetAllAreas(limit int, paginate bool, page int, filters map[string]string, userRole string, territoryID int) ([]response.AreaResponse, int64, error) {
	return s.repo.GetAllAreas(limit, paginate, page, filters, userRole, territoryID)
}

func (s *areaService) GetAreaByID(id int) (*response.AreaResponse, error) {
	return s.repo.GetAreaByID(id)
}

func (s *areaService) GetAreaByName(name string) (*response.AreaResponse, error) {
	return s.repo.GetAreaByName(name)
}

func (s *areaService) CreateArea(name *string) (*response.AreaResponse, error) {
	area := &models.Area{Name: *name}
	return s.repo.CreateArea(area)
}

func (s *areaService) UpdateArea(name *string, id int) (*response.AreaResponse, error) {
	area := &models.Area{Name: *name}
	return s.repo.UpdateArea(area, id)
}
