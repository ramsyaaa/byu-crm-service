package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/kpi-yae-range/validation"
)

type KpiYaeRangeService interface {
	GetKpiYaeRangeByDate(month uint, year uint) (*models.KpiYaeRange, error)
	CreateKpiYaeRange(kpiYaeRange *validation.CreateKpiYaeRangeRequest) (*models.KpiYaeRange, error)
}
