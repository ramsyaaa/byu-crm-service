package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/kpi-yae-range/repository"
	"byu-crm-service/modules/kpi-yae-range/validation"
	"encoding/json"
	"time"
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

func (s *kpiYaeRangeService) CreateKpiYaeRange(kpiYaeRange *validation.CreateKpiYaeRangeRequest) (*models.KpiYaeRange, error) {
	// Parse Dates from string to time.Time
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, kpiYaeRange.StartDate)
	if err != nil {
		return nil, err
	}
	endDate, err := time.Parse(layout, kpiYaeRange.EndDate)
	if err != nil {
		return nil, err
	}

	// Prepare data for JSON encoding
	var allTarget []map[string]string
	for i, target := range kpiYaeRange.Target {
		if target != "" {
			allTarget = append(allTarget, map[string]string{
				"name":   kpiYaeRange.Name[i],
				"target": target,
			})
		}
	}

	// Encode to JSON
	targetJSON, err := json.Marshal(allTarget)
	if err != nil {
		return nil, err
	}

	// Save model
	kpiYaeRangeModel := &models.KpiYaeRange{
		StartDate: &startDate,
		EndDate:   &endDate,
		Target:    string(targetJSON),
	}

	return s.repo.CreateKpiYaeRange(kpiYaeRangeModel)
}
