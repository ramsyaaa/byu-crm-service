package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/type/repository"
)

type typeService struct {
	repo repository.TypeRepository
}

func NewTypeService(repo repository.TypeRepository) TypeService {
	return &typeService{repo: repo}
}

func (s *typeService) GetAllTypes(limit int, paginate bool, page int, filters map[string]string, module string, category_name []string) ([]models.Type, int64, error) {
	return s.repo.GetAllTypes(limit, paginate, page, filters, module, category_name)
}
