package service

import "byu-crm-service/models"

type TypeService interface {
	GetAllTypes(limit int, paginate bool, page int, filters map[string]string, module string, category_name []string) ([]models.Type, int64, error)
	GetTypeByID(id int) (models.Type, error)
	GetTypeByNameAndModuleType(name string, moduleType string, categoryID int) (models.Type, error)
	CreateType(requestBody map[string]interface{}) (models.Type, error)
	UpdateType(id int, requestBody map[string]interface{}) (models.Type, error)
}
