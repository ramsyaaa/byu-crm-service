package service

type NotificationOneSignalService interface {
	SendNotification(requestBody map[string]string, rolesName []string, UserRole string, territoryID int, userID int) error
	AssignNotificationOneSignalToUsers(requestBody map[string]string, userIDs []int) error
	CreateSubscribeNotification(userID *uint, SubscribeID string, SubscribeType string) error
	DeleteSubscribeNotification(userID *uint, SubscribeID string) error
}
