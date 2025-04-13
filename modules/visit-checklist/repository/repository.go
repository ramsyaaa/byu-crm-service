package repository

import "byu-crm-service/models"

type VisitChecklistRepository interface {
	GetAllVisitChecklist() ([]models.VisitChecklist, error)
}
