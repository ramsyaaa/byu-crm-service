package service

import (
	"byu-crm-service/models"
	"time"
)

type PerformanceSkulIdService interface {
	ProcessPerformanceSkulId(data []string) error
	FindAll(limit, offset int, filters map[string]string, accountID int, page int, paginate bool) ([]models.PerformanceSkulId, int64, error)
	FindBySerialNumberMsisdn(serial string) (*models.PerformanceSkulId, error)
	FindByIdSkulId(idSkulId string) (*models.PerformanceSkulId, error)
	CreatePerformanceSkulID(account_id int, userName, idSkulId, msisdn string, registeredDate *time.Time, provider *string, batch *string, user_type *string) (*models.PerformanceSkulId, error)
}
