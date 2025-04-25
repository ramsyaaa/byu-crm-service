package service

import "byu-crm-service/models"

type TypeService interface {
	GetAllTypes(limit int, paginate bool, page int, filters map[string]string, module string, category_name []string) ([]models.Type, int64, error)
}
