package service

import "byu-crm-service/models"

type ConstantDataService interface {
	GetAllConstants(type_constant string, other_group string) ([]models.ConstantData, int64, error)
	CreateConstant(requestBody map[string]interface{}) (models.ConstantData, error)
	GetConstantByTypeAndValue(type_constant string, value string) (models.ConstantData, error)
}
