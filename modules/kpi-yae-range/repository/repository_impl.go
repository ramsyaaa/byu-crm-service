package repository

import (
	"byu-crm-service/models"
	"time"

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

	// Create a time.Time object for the first day of the month
	start := time.Date(int(year), time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	// Create a time.Time object for the last day of the month
	end := start.AddDate(0, 1, -1)

	// Convert to string in the format "YYYY-MM-DD"
	startDate := start.Format("2006-01-02")
	endDate := end.Format("2006-01-02")

	err := r.db.Where("start_date <= ? AND end_date >= ?", endDate, startDate).First(&kpi).Error
	if err != nil {
		return nil, err
	}

	return &kpi, nil
}

func (r *kpiYaeRangeRepository) CreateKpiYaeRange(kpiYaeRange *models.KpiYaeRange) (*models.KpiYaeRange, error) {
	err := r.db.Create(kpiYaeRange).Error
	if err != nil {
		return nil, err
	}

	return kpiYaeRange, nil
}
