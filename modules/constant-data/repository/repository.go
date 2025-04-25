package repository

import "byu-crm-service/models"

type ConstantDataRepository interface {
	GetAllConstants(limit int, paginate bool, page int, filters map[string]string, type_constant string) ([]models.ConstantData, int64, error)
}
