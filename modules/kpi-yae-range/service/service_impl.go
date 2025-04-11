package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/kpi-yae-range/repository"
)

type kpiYaeRangeService struct {
	repo repository.KpiYaeRangeRepository
}

func NewKpiYaeRangeService(repo repository.KpiYaeRangeRepository) KpiYaeRangeService {
	return &kpiYaeRangeService{repo: repo}
}

func (s *kpiYaeRangeService) GetKpiYaeRangeByDate(month uint, year uint) (*models.KpiYaeRange, error) {
	return s.repo.GetKpiYaeRangeByDate(month, year)
}
