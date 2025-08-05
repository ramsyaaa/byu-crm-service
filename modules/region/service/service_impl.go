package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/region/repository"
	"byu-crm-service/modules/region/response"
	"fmt"
)

type regionService struct {
	repo repository.RegionRepository
}

func NewRegionService(repo repository.RegionRepository) RegionService {
	return &regionService{repo: repo}
}

func (s *regionService) GetAllRegions(filters map[string]string, userRole string, territoryID int, withGeo bool) ([]response.RegionResponse, int64, error) {
	return s.repo.GetAllRegions(filters, userRole, territoryID, withGeo)
}

func (s *regionService) GetRegionByID(id int) (*response.RegionResponse, error) {
	return s.repo.GetRegionByID(id)
}

func (s *regionService) GetRegionByName(name string) (*response.RegionResponse, error) {
	return s.repo.GetRegionByName(name)
}

func (s *regionService) CreateRegion(name *string, area_id int) (*response.RegionResponse, error) {
	areaIDStr := fmt.Sprintf("%d", area_id)
	region := &models.Region{Name: *name, AreaID: &areaIDStr}
	return s.repo.CreateRegion(region)
}

func (s *regionService) UpdateRegion(name *string, area_id int, id int) (*response.RegionResponse, error) {
	areaIDStr := fmt.Sprintf("%d", area_id)
	region := &models.Region{Name: *name, AreaID: &areaIDStr}
	return s.repo.UpdateRegion(region, id)
}
