package repository

import (
	"byu-crm-service/models"
)

type PerformanceSkulIdRepository interface {
	Create(performanceSkulId *models.PerformanceSkulId) error
	FindBySerialNumberMsisdn(serial string) (*models.PerformanceSkulId, error)
	FindByIdSkulId(idSkulId string) (*models.PerformanceSkulId, error)
	Update(performance *models.PerformanceSkulId) error
}
