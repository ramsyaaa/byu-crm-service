package repository

import (
	"byu-crm-service/models"
	"errors"
	"fmt"

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

func (r *performanceIndianaRepository) GetDataInByUserAndMonth(
	userID int,
	month uint,
	year uint,
) (int, error) {

	monthStr := fmt.Sprintf("%04d-%02d", year, month)

	var value int

	err := r.db.
		Table("performance_indianas").
		Select("COALESCE(data_in, 0)").
		Where("user_id = ?", userID).
		Where("month = ?", monthStr).
		Limit(1).
		Scan(&value).Error

	if err != nil {
		return 0, err
	}

	// kalau record tidak ada → Scan tetap 0
	return value, nil
}
