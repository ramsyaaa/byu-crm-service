package repository

import (
	"byu-crm-service/models"
)

type PerformanceIndianaRepository interface {
	FindByUserAndMonth(userID int, month string) (*models.PerformanceIndiana, error)
	Create(data *models.PerformanceIndiana) error
	Update(data *models.PerformanceIndiana) error
}
