package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/type/repository"
	"strconv"
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

func (s *typeService) GetTypeByID(id int) (models.Type, error) {
	return s.repo.GetTypeByID(id)
}

func (s *typeService) GetTypeByNameAndModuleType(name string, moduleType string, categoryID int) (models.Type, error) {
	return s.repo.GetTypeByNameAndModuleType(name, moduleType, categoryID)
}

func (s *typeService) CreateType(requestBody map[string]interface{}) (models.Type, error) {
	var newType models.Type

	if val, ok := requestBody["category_id"].(string); ok && val != "" {
		if parsedVal, err := strconv.ParseUint(val, 10, 64); err == nil {
			temp := uint(parsedVal)
			newType.CategoryID = &temp
		}
	}

	if val, ok := requestBody["module_type"].(string); ok && val != "" {
		newType.ModuleType = &val
	}

	if val, ok := requestBody["name"].(string); ok && val != "" {
		newType.Name = &val
	}

	if val, ok := requestBody["description"].(string); ok && val != "" {
		newType.Description = &val
	}

	return s.repo.CreateType(newType)
}

func (s *typeService) UpdateType(id int, requestBody map[string]interface{}) (models.Type, error) {
	var typeData models.Type

	if val, ok := requestBody["category_id"].(string); ok && val != "" {
		if parsedVal, err := strconv.ParseUint(val, 10, 64); err == nil {
			temp := uint(parsedVal)
			typeData.CategoryID = &temp
		}
	}

	if val, ok := requestBody["module_type"].(string); ok && val != "" {
		typeData.ModuleType = &val
	}

	if val, ok := requestBody["name"].(string); ok && val != "" {
		typeData.Name = &val
	}

	if val, ok := requestBody["description"].(string); ok && val != "" {
		typeData.Description = &val
	}

	return s.repo.UpdateType(id, typeData)
}
