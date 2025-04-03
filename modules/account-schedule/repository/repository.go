package repository

import "byu-crm-service/models"

type AccountScheduleRepository interface {
	GetBySubject(subject_type string, subject_id uint) ([]models.AccountSchedule, error)
	DeleteBySubject(subject_type string, subject_id uint) error
	Insert(accountSchedule []models.AccountSchedule) error
}
