package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/category/repository"
)

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) GetAllCategories(limit int, paginate bool, page int, filters map[string]string, module string) ([]models.Category, int64, error) {
	return s.repo.GetAllCategories(limit, paginate, page, filters, module)
}

func (s *categoryService) GetCategoryByID(id int) (*models.Category, error) {
	return s.repo.GetCategoryByID(id)
}

func (s *categoryService) GetCategoryByNameAndModuleType(name string, moduleType string) (*models.Category, error) {
	return s.repo.GetCategoryByNameAndModuleType(name, moduleType)
}

func (s *categoryService) CreateCategory(requestBody map[string]interface{}) (*models.Category, error) {
	var newCategory models.Category

	if val, ok := requestBody["module_type"].(string); ok && val != "" {
		newCategory.ModuleType = &val
	}

	if val, ok := requestBody["name"].(string); ok && val != "" {
		newCategory.Name = &val
	}

	return s.repo.CreateCategory(newCategory)
}

func (s *categoryService) UpdateCategory(id int, requestBody map[string]interface{}) (models.Category, error) {
	var category models.Category

	if val, ok := requestBody["module_type"].(string); ok && val != "" {
		category.ModuleType = &val
	}

	if val, ok := requestBody["name"].(string); ok && val != "" {
		category.Name = &val
	}

	return s.repo.UpdateCategory(id, category)
}
