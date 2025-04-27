package repository

import "byu-crm-service/models"

type ConstantDataRepository interface {
	GetAllConstants(limit int, paginate bool, page int, filters map[string]string, type_constant string, other_group string) ([]models.ConstantData, int64, error)
	CreateConstant(constantData models.ConstantData) (models.ConstantData, error)
	GetConstantByTypeAndValue(type_constant string, value string) (models.ConstantData, error)
}
