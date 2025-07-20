package service

import (
	"byu-crm-service/modules/notification/repository"
	"byu-crm-service/modules/notification/response"
	userRepo "byu-crm-service/modules/user/repository"
)

type notificationService struct {
	repo     repository.NotificationRepository
	userRepo userRepo.UserRepository
}

func NewNotificationService(repo repository.NotificationRepository, userRepo userRepo.UserRepository) NotificationService {
	return &notificationService{repo: repo, userRepo: userRepo}
}

func (s *notificationService) GetAllNotifications(filters map[string]string, limit, userID int) ([]response.NotificationResponse, int64, error) {
	return s.repo.GetAllNotifications(filters, limit, userID)
}

func (s *notificationService) AssignNotificationToUsers(requestBody map[string]string, userIDs []int) error {
	return s.repo.CreateNotification(requestBody, userIDs)
}

func (s *notificationService) CreateNotification(requestBody map[string]string, rolesName []string, userRole string, territoryID int, userID int) error {
	if userID != 0 {
		userIDs := []int{userID}
		return s.repo.CreateNotification(requestBody, userIDs)
	} else {
		filters := map[string]string{
			"search":     "",
			"order_by":   "id",
			"order":      "DESC",
			"start_date": "",
			"end_date":   "",
		}
		users, _, err := s.userRepo.GetAllUsers(0, false, 1, filters, rolesName, false, userRole, territoryID)

		if err != nil {
			return err
		}

		var userIDs []int
		for _, u := range users {
			userIDs = append(userIDs, int(u.ID))
		}

		return s.repo.CreateNotification(requestBody, userIDs)
	}

}

func (s *notificationService) GetByNotificationId(notificationID uint, userID uint) (*response.NotificationResponse, error) {
	return s.repo.GetByNotificationId(notificationID, userID)
}

func (s *notificationService) MarkNotificationAsRead(notificationID uint, userID uint) error {
	return s.repo.MarkNotificationAsRead(notificationID, userID)
}

func (s *notificationService) MarkAllNotificationsAsRead(userID int) error {
	return s.repo.MarkAllNotificationsAsRead(userID)
}
