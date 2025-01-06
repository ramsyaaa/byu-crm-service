package repository

import (
	"byu-crm-service/models"
)

type PerformanceDigiposRepository interface {
	Create(performanceDigipos *models.PerformanceDigipos) error
}
