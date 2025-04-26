package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/area/repository"
)

type areaService struct {
	repo repository.AreaRepository
}

func NewAreaService(repo repository.AreaRepository) AreaService {
	return &areaService{repo: repo}
}

func (s *areaService) GetAllAreas(limit int, paginate bool, page int, filters map[string]string) ([]models.Area, int64, error) {
	return s.repo.GetAllAreas(limit, paginate, page, filters)
}

func (s *areaService) GetAreaByID(id int) (*models.Area, error) {
	return s.repo.GetAreaByID(id)
}

func (s *areaService) GetAreaByName(name string) (*models.Area, error) {
	return s.repo.GetAreaByName(name)
}

func (s *areaService) CreateArea(name *string) (*models.Area, error) {
	area := &models.Area{Name: *name}
	return s.repo.CreateArea(area)
}

func (s *areaService) UpdateArea(name *string, id int) (*models.Area, error) {
	area := &models.Area{Name: *name}
	return s.repo.UpdateArea(area, id)
}
