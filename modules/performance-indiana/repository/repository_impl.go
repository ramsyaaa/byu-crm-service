package repository

import (
	"byu-crm-service/models"
	"errors"

	"gorm.io/gorm"
)

type performanceIndianaRepository struct {
	db *gorm.DB
}

func NewPerformanceIndianaRepository(db *gorm.DB) PerformanceIndianaRepository {
	return &performanceIndianaRepository{db: db}
}

func (r *performanceIndianaRepository) FindByUserAndMonth(
	userID int,
	month string,
) (*models.PerformanceIndiana, error) {

	var data models.PerformanceIndiana

	err := r.db.
		Where("user_id = ? AND month = ?", userID, month).
		First(&data).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &data, nil
}

func (r *performanceIndianaRepository) Create(
	data *models.PerformanceIndiana,
) error {
	return r.db.Create(data).Error
}

func (r *performanceIndianaRepository) Update(
	data *models.PerformanceIndiana,
) error {
	return r.db.Save(data).Error
}
