package repository

import (
	"byu-crm-service/modules/notification/response"
)

type NotificationRepository interface {
	GetAllNotifications(filters map[string]string, limit, userID int) ([]response.NotificationResponse, int64, error)
	CreateNotification(requestBody map[string]string, userIDs []int) error
	GetByNotificationId(notificationID uint, userID uint) (*response.NotificationResponse, error)
}
