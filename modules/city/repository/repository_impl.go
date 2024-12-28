package repository

import (
	"byu-crm-service/models"
	"errors"

	"gorm.io/gorm"
)

type cityRepository struct {
	db *gorm.DB
}

func NewCityRepository(db *gorm.DB) CityRepository {
	return &cityRepository{db: db}
}

func (r *cityRepository) FindByID(id uint) (*models.City, error) {
	var city models.City
	if err := r.db.First(&city, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not found is not an error
		}
		return nil, err
	}
	return &city, nil
}

func (r *cityRepository) FindByName(name string) (*models.City, error) {
	var city models.City
	if err := r.db.Where("name = ?", name).First(&city).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not found is not an error
		}
		return nil, err
	}
	return &city, nil
}
