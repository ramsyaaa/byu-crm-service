package repository

import "byu-crm-service/models"

type NotificationOneSignalRepository interface {
	SendNotification(requestBody map[string]string, playerID []string) error
	GetSubscribeNotificationsByUserIDs(userIDs []int) ([]models.SubscribeNotification, error)
	CreateSubscribeNotification(dataSubscribe *models.SubscribeNotification) error
	DeleteSubscribeNotification(userID *uint, subscriptionID string) error
}
