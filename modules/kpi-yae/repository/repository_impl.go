package repository

import (
	"byu-crm-service/models"

	"gorm.io/gorm"
)

type kpiYaeRepository struct {
	db *gorm.DB
}

func NewKpiYaeRepository(db *gorm.DB) KpiYaeRepository {
	return &kpiYaeRepository{db: db}
}

func (r *kpiYaeRepository) GetKpiYaeByName(name string) (*models.KpiYae, error) {
	var kpi models.KpiYae

	err := r.db.Where("BINARY name = ?", name).First(&kpi).Error

	if err != nil {
		return nil, err
	}

	return &kpi, nil
}
