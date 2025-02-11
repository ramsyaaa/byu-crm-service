package repository

import (
	"byu-crm-service/models"
)

type PerformanceDigiposRepository interface {
	Create(performanceDigipos *models.PerformanceDigipos) error
	FindByIdImport(idImport string) (*models.PerformanceDigipos, error)
	Update(performance *models.PerformanceDigipos) error
}
