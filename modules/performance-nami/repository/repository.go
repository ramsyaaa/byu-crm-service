package repository

import (
	"byu-crm-service/models"
)

type PerformanceNamiRepository interface {
	Create(performanceNami *models.PerformanceNami) error
}
