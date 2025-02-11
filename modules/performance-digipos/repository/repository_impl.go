package repository

import (
	"byu-crm-service/models"
	"errors"

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

func (r *performanceDigiposRepository) FindByIdImport(idImport string) (*models.PerformanceDigipos, error) {
	var performance models.PerformanceDigipos
	err := r.db.Where("id_import = ?", idImport).First(&performance).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Data tidak ditemukan, kembalikan nil
		}
		return nil, err
	}
	return &performance, nil
}

func (r *performanceDigiposRepository) Update(performance *models.PerformanceDigipos) error {
	return r.db.Save(performance).Error
}
