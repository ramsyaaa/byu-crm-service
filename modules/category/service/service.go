package service

import "byu-crm-service/models"

type CategoryService interface {
	GetAllCategories(limit int, paginate bool, page int, filters map[string]string, module string) ([]models.Category, int64, error)
	GetCategoryByID(id int) (*models.Category, error)
	GetCategoryByName(name string) (*models.Category, error)
	GetCategoryByNameAndModuleType(name string, moduleType string) (*models.Category, error)
	CreateCategory(requestBody map[string]interface{}) (*models.Category, error)
	UpdateCategory(id int, requestBody map[string]interface{}) (models.Category, error)
}
