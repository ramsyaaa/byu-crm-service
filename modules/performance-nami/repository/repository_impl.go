package repository

import (
	"byu-crm-service/models"
	"errors"

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

func (r *performanceNamiRepository) FindBySerialNumberMsisdn(serial string) (*models.PerformanceNami, error) {
	var performanceNami models.PerformanceNami
	err := r.db.Where("serial_number_msisdn = ?", serial).First(&performanceNami).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &performanceNami, nil
}
