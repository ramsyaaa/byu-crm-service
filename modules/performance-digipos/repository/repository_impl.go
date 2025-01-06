package repository

import (
	"byu-crm-service/models"

	"gorm.io/gorm"
)

type performanceDigiposRepository struct {
	db *gorm.DB
}

func NewPerformanceDigiposRepository(db *gorm.DB) PerformanceDigiposRepository {
	return &performanceDigiposRepository{db: db}
}

func (r *performanceDigiposRepository) Create(performanceDigipos *models.PerformanceDigipos) error {
	return r.db.Create(performanceDigipos).Error
}
