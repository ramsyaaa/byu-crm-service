package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/visit-checklist/repository"
)

type visitChecklistService struct {
	repo repository.VisitChecklistRepository
}

func NewVisitChecklistService(repo repository.VisitChecklistRepository) VisitChecklistService {
	return &visitChecklistService{repo: repo}
}

func (s *visitChecklistService) GetAllVisitChecklist() ([]models.VisitChecklist, error) {
	return s.repo.GetAllVisitChecklist()
}
