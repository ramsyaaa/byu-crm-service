package service

import "byu-crm-service/models"

type AccountService interface {
	GetAllAccounts(limit, page int, search, userRole, territoryID string) ([]models.Account, map[string]interface{}, error)
	ProcessAccount(data []string) error
}
