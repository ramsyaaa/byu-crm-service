package service

import "byu-crm-service/models"

type ConstantDataService interface {
	GetAllConstants(limit int, paginate bool, page int, filters map[string]string, type_constant string) ([]models.ConstantData, int64, error)
}
