package repository

import (
	"byu-crm-service/models"

	"gorm.io/gorm"
)

type performanceNamiRepository struct {
	db *gorm.DB
}

func NewPerformanceNamiRepository(db *gorm.DB) PerformanceNamiRepository {
	return &performanceNamiRepository{db: db}
}

func (r *performanceNamiRepository) Create(performanceNami *models.PerformanceNami) error {
	return r.db.Create(performanceNami).Error
}
