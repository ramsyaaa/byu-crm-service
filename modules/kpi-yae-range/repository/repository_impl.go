package repository

import (
	"byu-crm-service/models"
	"fmt"

	"gorm.io/gorm"
)

type kpiYaeRangeRepository struct {
	db *gorm.DB
}

func NewKpiYaeRangeRepository(db *gorm.DB) KpiYaeRangeRepository {
	return &kpiYaeRangeRepository{db: db}
}

func (r *kpiYaeRangeRepository) GetKpiYaeRangeByDate(month uint, year uint) (*models.KpiYaeRange, error) {
	var kpi models.KpiYaeRange

	// Buat batas awal dan akhir bulan
	startDate := fmt.Sprintf("%04d-%02d-01", year, month)
	endDate := fmt.Sprintf("%04d-%02d-31", year, month)

	err := r.db.Where("start_date <= ? AND end_date >= ?", endDate, startDate).
		First(&kpi).Error

	if err != nil {
		return nil, err
	}

	return &kpi, nil
}
