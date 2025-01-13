package repository

import (
	"byu-crm-service/models"
)

type PerformanceSkulIdRepository interface {
	Create(performanceSkulId *models.PerformanceSkulId) error
	FindBySerialNumberMsisdn(serial string) (*models.PerformanceSkulId, error)
}
