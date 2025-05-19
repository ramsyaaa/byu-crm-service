package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/visit-history/repository"
	"encoding/json"
)

type visitHistoryService struct {
	repo repository.VisitHistoryRepository
}

func NewVisitHistoryService(repo repository.VisitHistoryRepository) VisitHistoryService {
	return &visitHistoryService{repo: repo}
}

func (s *visitHistoryService) CreateVisitHistory(user_id int, subject_type string, subject_id int, absence_user_id int, kpiYae map[string]*int, description *string, detailVisit map[string]string) (*models.VisitHistory, error) {
	convertedUserID := uint(user_id)
	kpiJSON, err := json.Marshal(kpiYae)
	if err != nil {
		return nil, err
	}
	kpiJSONString := string(kpiJSON)
	visitDetail, err := json.Marshal(detailVisit)
	if err != nil {
		return nil, err
	}
	visitDetailString := string(visitDetail)
	VisitHistory := &models.VisitHistory{
		UserID:        &convertedUserID,
		SubjectType:   &subject_type,
		SubjectID:     func(v int) *uint { u := uint(v); return &u }(subject_id),
		AbsenceUserID: func(v int) *uint { u := uint(v); return &u }(absence_user_id),
		Target:        &kpiJSONString,
		DetailVisit:   &visitDetailString,
		Description:   *description,
	}
	return s.repo.CreateVisitHistory(VisitHistory)
}

func (s *visitHistoryService) CountVisitHistory(user_id int, month uint, year uint, kpi_name string) (int, error) {
	return s.repo.CountVisitHistory(user_id, month, year, kpi_name)
}
