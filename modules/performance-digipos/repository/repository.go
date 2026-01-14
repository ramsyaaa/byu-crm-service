package repository

import (
	"byu-crm-service/models"
)

type PerformanceDigiposRepository interface {
	Create(performanceDigipos *models.PerformanceDigipos) error
	FindByIdImport(idImport string) (*models.PerformanceDigipos, error)
	Update(performance *models.PerformanceDigipos) error
	CountPerformanceByUserYaeCode(userID int, month uint, year uint) (int, error)
}
