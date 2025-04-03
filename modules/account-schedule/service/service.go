package service

import "byu-crm-service/models"

type AccountScheduleService interface {
	GetBySubject(subject_type string, subject_id uint) ([]models.AccountSchedule, error)
	Insert(requestBody map[string]interface{}, subject_type string, subject_id uint) ([]models.AccountSchedule, error)
}
