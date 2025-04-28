package repository

import "byu-crm-service/models"

type CategoryRepository interface {
	GetAllCategories(limit int, paginate bool, page int, filters map[string]string, module string) ([]models.Category, int64, error)
	GetCategoryByID(id int) (*models.Category, error)
	GetCategoryByName(name string) (*models.Category, error)
	GetCategoryByNameAndModuleType(name string, moduleType string) (*models.Category, error)
	CreateCategory(requestBody models.Category) (*models.Category, error)
	UpdateCategory(id int, updateCategory models.Category) (models.Category, error)
}
