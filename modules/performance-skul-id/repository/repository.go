package repository

import (
	"byu-crm-service/models"
)

type PerformanceSkulIdRepository interface {
	Create(performanceSkulId *models.PerformanceSkulId) error
	FindBySerialNumberMsisdn(serial string) (*models.PerformanceSkulId, error)
	FindByIdSkulId(idSkulId string) (*models.PerformanceSkulId, error)
	Update(performance *models.PerformanceSkulId) error
	UpdateByFields(id uint, fields map[string]interface{}) error
	FindAll(limit, offset int, filters map[string]string, accountID int, page int, paginate bool) ([]models.PerformanceSkulId, int64, error)
}
