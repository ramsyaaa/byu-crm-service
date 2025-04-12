package repository

import "byu-crm-service/models"

type KpiYaeRangeRepository interface {
	GetKpiYaeRangeByDate(month uint, year uint) (*models.KpiYaeRange, error)
	CreateKpiYaeRange(kpiYaeRange *models.KpiYaeRange) (*models.KpiYaeRange, error)
}
