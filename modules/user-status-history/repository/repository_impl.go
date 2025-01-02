package repository

import (
	"byu-crm-service/models"

	"gorm.io/gorm"
)

type userStatusHistoryRepository struct {
	db *gorm.DB
}

func NewUserStatusHistoryRepository(db *gorm.DB) *userStatusHistoryRepository {
	return &userStatusHistoryRepository{db: db}
}

func (r *userStatusHistoryRepository) FindAllByStatus(status string) ([]models.UserStatusHistory, error) {
	var userStatusHistories []models.UserStatusHistory
	if err := r.db.Where("status = ?", status).Find(&userStatusHistories).Error; err != nil {
		return nil, err
	}
	return userStatusHistories, nil
}

// Update updates a single user status history record.
func (r *userStatusHistoryRepository) Update(userStatusHistory *models.UserStatusHistory) error {
	return r.db.Save(userStatusHistory).Error
}
