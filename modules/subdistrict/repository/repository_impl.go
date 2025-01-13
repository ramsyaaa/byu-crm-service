package repository

import (
	"byu-crm-service/models"
	"errors"

	"gorm.io/gorm"
)

type subdistrictRepository struct {
	db *gorm.DB
}

func NewSubdistrictRepository(db *gorm.DB) SubdistrictRepository {
	return &subdistrictRepository{db: db}
}

func (r *subdistrictRepository) FindByID(id uint) (*models.Subdistrict, error) {
	var Subdistrict models.Subdistrict
	if err := r.db.First(&Subdistrict, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not found is not an error
		}
		return nil, err
	}
	return &Subdistrict, nil
}

func (r *subdistrictRepository) FindByName(name string) (*models.Subdistrict, error) {
	var Subdistrict models.Subdistrict
	if err := r.db.Where("name = ?", name).First(&Subdistrict).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not found is not an error
		}
		return nil, err
	}
	return &Subdistrict, nil
}
