package service

import "byu-crm-service/models"

type KpiYaeRangeService interface {
	GetKpiYaeRangeByDate(month uint, year uint) (*models.KpiYaeRange, error)
}
