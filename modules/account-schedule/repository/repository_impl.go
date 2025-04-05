package repository

import (
	"byu-crm-service/models"

	"gorm.io/gorm"
)

type accountScheduleRepository struct {
	db *gorm.DB
}

func NewAccountScheduleRepository(db *gorm.DB) AccountScheduleRepository {
	return &accountScheduleRepository{db: db}
}

func (r *accountScheduleRepository) GetBySubject(subject_type string, subject_id uint) ([]models.AccountSchedule, error) {
	var accountSchedule []models.AccountSchedule

	if err := r.db.Where("subject_type = ? AND subject_id = ?", subject_type, subject_id).First(&accountSchedule).Error; err != nil {
		return nil, err
	}

	return accountSchedule, nil
}

func (r *accountScheduleRepository) DeleteBySubject(subject_type string, subject_id uint) error {
	return r.db.Where("subject_type = ? AND subject_id = ?", subject_type, subject_id).
		Delete(&models.AccountSchedule{}).Error
}

func (r *accountScheduleRepository) Insert(accountSchedule []models.AccountSchedule) error {
	return r.db.Create(&accountSchedule).Error
}
