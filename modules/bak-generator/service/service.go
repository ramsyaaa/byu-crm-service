package service

import "byu-crm-service/models"

type BakGeneratorService interface {
	CreateBak(reqMap map[string]interface{}, user_id uint) error
	GetBakByID(id uint) (*models.BakFile, error)
	GetAllBak(limit int, paginate bool, page int, filters map[string]string) ([]models.BakFile, int, error)
}
