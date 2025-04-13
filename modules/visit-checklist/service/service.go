package service

import "byu-crm-service/models"

type VisitChecklistService interface {
	GetAllVisitChecklist() ([]models.VisitChecklist, error)
}
