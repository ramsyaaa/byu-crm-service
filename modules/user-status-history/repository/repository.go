package repository

import "byu-crm-service/models"

type UserStatusHistoryRepository interface {
	FindAllByStatus(status string) ([]models.UserStatusHistory, error)
	Update(userStatusHistory *models.UserStatusHistory) error
}
