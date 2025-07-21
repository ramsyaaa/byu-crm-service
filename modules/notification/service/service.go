package service

import (
	"byu-crm-service/modules/notification/response"
)

type NotificationService interface {
	GetAllNotifications(filters map[string]string, limit, userID int) ([]response.NotificationResponse, int64, error)
	CreateNotification(requestBody map[string]string, rolesName []string, UserRole string, territoryID int, userID int) error
	AssignNotificationToUsers(requestBody map[string]string, userIDs []int) error
	GetByNotificationId(notificationID uint, userID uint) (*response.NotificationResponse, error)
	MarkNotificationAsRead(notificationID uint, userID uint) error
	MarkAllNotificationsAsRead(userID int) error
	MarkNotificationAsReadBySubjectID(subjectType string, subjectID uint, userID int) error
}
