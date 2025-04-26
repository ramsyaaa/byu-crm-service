package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/subdistrict/repository"
	"byu-crm-service/modules/subdistrict/response"
	"fmt"
)

type subdistrictService struct {
	repo repository.SubdistrictRepository
}

func NewSubdistrictService(repo repository.SubdistrictRepository) SubdistrictService {
	return &subdistrictService{repo: repo}
}

func (s *subdistrictService) GetAllSubdistricts(limit int, paginate bool, page int, filters map[string]string, userRole string, territoryID int) ([]response.SubdistrictResponse, int64, error) {
	return s.repo.GetAllSubdistricts(limit, paginate, page, filters, userRole, territoryID)
}

func (s *subdistrictService) GetSubdistrictByID(id int) (*response.SubdistrictResponse, error) {
	return s.repo.GetSubdistrictByID(id)
}

func (s *subdistrictService) GetSubdistrictByName(name string) (*response.SubdistrictResponse, error) {
	return s.repo.GetSubdistrictByName(name)
}

func (s *subdistrictService) CreateSubdistrict(name *string, city_id int) (*response.SubdistrictResponse, error) {
	cityIDStr := fmt.Sprintf("%d", city_id)
	subdistrict := &models.Subdistrict{Name: *name, CityID: &cityIDStr}
	return s.repo.CreateSubdistrict(subdistrict)
}

func (s *subdistrictService) UpdateSubdistrict(name *string, city_id int, id int) (*response.SubdistrictResponse, error) {
	cityIDStr := fmt.Sprintf("%d", city_id)
	subdistrict := &models.Subdistrict{Name: *name, CityID: &cityIDStr}
	return s.repo.UpdateSubdistrict(subdistrict, id)
}
