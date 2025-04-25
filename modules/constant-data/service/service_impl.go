package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/constant-data/repository"
)

type constantDataService struct {
	repo repository.ConstantDataRepository
}

func NewConstantDataService(repo repository.ConstantDataRepository) ConstantDataService {
	return &constantDataService{repo: repo}
}

func (s *constantDataService) GetAllConstants(limit int, paginate bool, page int, filters map[string]string, type_constant string) ([]models.ConstantData, int64, error) {
	return s.repo.GetAllConstants(limit, paginate, page, filters, type_constant)
}
