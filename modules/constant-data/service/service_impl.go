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

func (s *constantDataService) GetAllConstants(limit int, paginate bool, page int, filters map[string]string, type_constant string, other_group string) ([]models.ConstantData, int64, error) {
	return s.repo.GetAllConstants(limit, paginate, page, filters, type_constant, other_group)
}

func (s *constantDataService) CreateConstant(requestBody map[string]interface{}) (models.ConstantData, error) {
	var constantData models.ConstantData

	if val, ok := requestBody["type"].(string); ok && val != "" {
		constantData.Type = &val
	}
	if val, ok := requestBody["value"].(string); ok && val != "" {
		constantData.Value = &val
	}
	if val, ok := requestBody["label"].(string); ok && val != "" {
		constantData.Label = &val
	}
	if val, ok := requestBody["other_group"].(string); ok && val != "" {
		constantData.OtherGroup = &val
	}

	return s.repo.CreateConstant(constantData)
}

func (s *constantDataService) GetConstantByTypeAndValue(type_constant string, value string) (models.ConstantData, error) {
	return s.repo.GetConstantByTypeAndValue(type_constant, value)
}
