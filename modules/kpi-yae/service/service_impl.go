package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/kpi-yae/repository"
)

type kpiYaeService struct {
	repo repository.KpiYaeRepository
}

func NewKpiYaeService(repo repository.KpiYaeRepository) KpiYaeService {
	return &kpiYaeService{repo: repo}
}

func (s *kpiYaeService) GetKpiYaeByName(name string) (*models.KpiYae, error) {
	return s.repo.GetKpiYaeByName(name)
}
