package service

import (
	"byu-crm-service/modules/user-status-history/repository"
	"time"
)

type userStatusHistoryService struct {
	repo repository.UserStatusHistoryRepository
}

func NewUserStatusHistoryService(repo repository.UserStatusHistoryRepository) UserStatusHistoryService {
	return &userStatusHistoryService{repo: repo}
}

func (s *userStatusHistoryService) UpdateDate() error {
	// Find all records with status 'active'
	activeStatuses, err := s.repo.FindAllByStatus("active")
	if err != nil {
		return err
	}
	if len(activeStatuses) == 0 {
		return nil // No active statuses found
	}

	// Calculate the end of the current month
	now := time.Now()
	endOfMonth := time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 0, now.Location())

	// Update end_date for each active status
	for i := range activeStatuses {
		activeStatuses[i].EndDate = endOfMonth
		if err := s.repo.Update(&activeStatuses[i]); err != nil {
			return err
		}
	}

	return nil
}
