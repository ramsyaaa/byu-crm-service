package service

import "byu-crm-service/models"

type VisitHistoryService interface {
	CreateVisitHistory(user_id int, subject_type string, subject_id int, absence_user_id int, kpiYae map[string]*int, description *string) (*models.VisitHistory, error)
	CountVisitHistory(user_id int, month uint, year uint, kpi_name string) (int, error)
}
