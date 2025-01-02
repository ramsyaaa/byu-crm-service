package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/subdistrict/repository"
)

type subdistrictService struct {
	repo repository.SubdistrictRepository
}

func NewSubdistrictService(repo repository.SubdistrictRepository) SubdistrictService {
	return &subdistrictService{repo: repo}
}

func (s *subdistrictService) GetSubdistrictByID(id uint) (*models.Subdistrict, error) {
	return s.repo.FindByID(id)
}

func (s *subdistrictService) GetSubdistrictByName(name string) (*models.Subdistrict, error) {
	return s.repo.FindByName(name)
}
