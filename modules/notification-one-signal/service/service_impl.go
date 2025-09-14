package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/notification-one-signal/repository"
	userRepo "byu-crm-service/modules/user/repository"
)

type notificationService struct {
	repo     repository.NotificationOneSignalRepository
	userRepo userRepo.UserRepository
}

func NewNotificationOneSignalService(repo repository.NotificationOneSignalRepository, userRepo userRepo.UserRepository) NotificationOneSignalService {
	return &notificationService{repo: repo, userRepo: userRepo}
}

func (s *notificationService) AssignNotificationOneSignalToUsers(requestBody map[string]string, userIDs []int) error {
	subscriptions, err := s.repo.GetSubscribeNotificationsByUserIDs(userIDs)

	if err != nil {
		return err
	}

	if len(subscriptions) == 0 {
		return nil
	}

	var playerIDs []string
	for _, sub := range subscriptions {
		playerIDs = append(playerIDs, sub.SubscribeID)
	}

	return s.repo.SendNotification(requestBody, playerIDs)
}

func (s *notificationService) SendNotification(requestBody map[string]string, rolesName []string, userRole string, territoryID int, userID int) error {
	if userID != 0 {
		subscriptions, err := s.repo.GetSubscribeNotificationsByUserIDs([]int{userID})
		if err != nil {
			return err
		}
		if len(subscriptions) == 0 {
			return nil
		}
		var playerIDs []string
		for _, sub := range subscriptions {
			playerIDs = append(playerIDs, sub.SubscribeID)
		}
		return s.repo.SendNotification(requestBody, playerIDs)
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

		subscription, err := s.repo.GetSubscribeNotificationsByUserIDs(userIDs)
		if err != nil {
			return err
		}
		if len(subscription) == 0 {
			return nil
		}
		var playerIDs []string
		for _, sub := range subscription {
			playerIDs = append(playerIDs, sub.SubscribeID)
		}

		return s.repo.SendNotification(requestBody, playerIDs)
	}

}

func (s *notificationService) CreateSubscribeNotification(userID *uint, SubscribeID string, SubscribeType string) error {
	userData := &models.SubscribeNotification{
		UserID:        userID,
		SubscribeID:   SubscribeID,
		SubscribeType: SubscribeType,
	}

	existingSub, _ := s.repo.GetSubscribeNotificationBySubscriptionID(SubscribeID)

	if existingSub != nil {
		return s.repo.UpdateSubscribeNotificationBySubscribeID(SubscribeID, userData)
	} else {
		return s.repo.CreateSubscribeNotification(userData)
	}
}

func (s *notificationService) DeleteSubscribeNotification(userID *uint, SubscribeID string) error {
	return s.repo.DeleteSubscribeNotification(userID, SubscribeID)
}
