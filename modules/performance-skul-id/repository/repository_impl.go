package repository

import (
	"byu-crm-service/models"
	"errors"

	"gorm.io/gorm"
)

type performanceSkulIdRepository struct {
	db *gorm.DB
}

func NewPerformanceSkulIdRepository(db *gorm.DB) PerformanceSkulIdRepository {
	return &performanceSkulIdRepository{db: db}
}

func (r *performanceSkulIdRepository) Create(performanceSkulId *models.PerformanceSkulId) error {
	return r.db.Create(performanceSkulId).Error
}

func (r *performanceSkulIdRepository) FindBySerialNumberMsisdn(serial string) (*models.PerformanceSkulId, error) {
	var performanceSkulId models.PerformanceSkulId
	err := r.db.Where("serial_number_msisdn = ?", serial).First(&performanceSkulId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &performanceSkulId, nil
}
