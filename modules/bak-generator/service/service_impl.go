package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/bak-generator/repository"
)

type bakGeneratorService struct {
	repo repository.BakGeneratorRepository
}

func NewBakGeneratorService(repo repository.BakGeneratorRepository) BakGeneratorService {
	return &bakGeneratorService{repo: repo}
}

func (s *bakGeneratorService) CreateBak(reqMap map[string]interface{}, user_id uint) error {
	return s.repo.CreateBak(reqMap, user_id)
}

func (s *bakGeneratorService) GetBakByID(id uint) (*models.BakFile, error) {
	return s.repo.GetBakByID(id)
}

func (s *bakGeneratorService) GetAllBak(limit int, paginate bool, page int, filters map[string]string) ([]models.BakFile, int, error) {
	return s.repo.GetAllBak(limit, paginate, page, filters)
}
