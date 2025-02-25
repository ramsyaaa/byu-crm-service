package repository

import "byu-crm-service/models"

type AccountRepository interface {
	FindByAccountCode(code string) (*models.Account, error)
	FindByAccountName(account_name string) (*models.Account, error)
	GetFilteredAccounts(limit, page int, search, userRole, territoryID string) ([]models.Account, int, error)
	Create(account *models.Account) error
	UpdateFields(id uint, fields map[string]interface{}) error
}
