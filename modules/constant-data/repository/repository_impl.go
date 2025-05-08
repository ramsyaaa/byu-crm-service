package repository

import (
	"byu-crm-service/models"

	"gorm.io/gorm"
)

type constantDataRepository struct {
	db *gorm.DB
}

func NewConstantDataRepository(db *gorm.DB) ConstantDataRepository {
	return &constantDataRepository{db: db}
}

func (r *constantDataRepository) GetAllConstants(type_constant string, other_group string) ([]models.ConstantData, int64, error) {
	var constant_data []models.ConstantData
	var total int64

	query := r.db.Model(&models.ConstantData{})

	// Apply type_constant filter
	if type_constant != "" {
		query = query.Where("constant_data.type = ?", type_constant)
	}

	// Apply other_group filter
	if other_group != "" {
		query = query.Where("constant_data.other_group = ?", other_group)
	}

	// Get total count before applying pagination
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Eksekusi query akhir
	err = query.Find(&constant_data).Error
	return constant_data, total, err
}

func (r *constantDataRepository) CreateConstant(constantData models.ConstantData) (models.ConstantData, error) {
	err := r.db.Create(&constantData).Error
	return constantData, err
}

func (r *constantDataRepository) GetConstantByTypeAndValue(type_constant string, value string) (models.ConstantData, error) {
	var constantData models.ConstantData
	err := r.db.Where("type = ? AND value = ?", type_constant, value).First(&constantData).Error
	if err != nil {
		return constantData, err
	}
	return constantData, nil
}
