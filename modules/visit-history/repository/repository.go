package repository

import "byu-crm-service/models"

type VisitHistoryRepository interface {
	CreateVisitHistory(VisitHistory *models.VisitHistory) (*models.VisitHistory, error)
	CountVisitHistory(user_id int, month uint, year uint, kpi_name string) (int, error)
}
