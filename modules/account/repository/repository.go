package repository

import "byu-crm-service/models"

type AccountRepository interface {
	FindByAccountCode(code string) (*models.Account, error)
	FindByAccountName(account_name string) (*models.Account, error)
	GetAllAccounts(limit int, paginate bool, page int, filters map[string]string, userRole string, territoryID int, userID int) ([]models.Account, int64, error)
	CreateAccount(requestBody map[string]string, userID int) ([]models.Account, error)
	UpdateAccount(requestBody map[string]string, accountID int, userID int) ([]models.Account, error)
	GetFilteredAccounts(limit, page int, search, userRole, territoryID string) ([]models.Account, int, error)
	Create(account *models.Account) error
	UpdateFields(id uint, fields map[string]interface{}) error
}
