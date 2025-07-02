package service

import (
	"byu-crm-service/modules/notification/response"
)

type NotificationService interface {
	GetAllNotifications(filters map[string]string, limit, userID int) ([]response.NotificationResponse, int64, error)
	CreateNotification(requestBody map[string]string, rolesName []string, territoryType string, territoryID int, userID int) error
}
