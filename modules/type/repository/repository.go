package repository

import "byu-crm-service/models"

type TypeRepository interface {
	GetAllTypes(limit int, paginate bool, page int, filters map[string]string, module string, category_name []string) ([]models.Type, int64, error)
	GetTypeByID(id int) (models.Type, error)
	GetTypeByName(name string) (models.Type, error)
	GetTypeByNameAndModuleType(name string, moduleType string, categoryID int) (models.Type, error)
	CreateType(newType models.Type) (models.Type, error)
	UpdateType(id int, updateType models.Type) (models.Type, error)
}
