package service

import "byu-crm-service/models"

type AccountService interface {
	GetAllAccounts(limit int, paginate bool, page int, filters map[string]string, userRole string, territoryID int, userID int, onlyUserPic bool) ([]models.Account, int64, error)
	CreateAccount(requestBody map[string]interface{}, userID int) ([]models.Account, error)
	UpdateAccount(requestBody map[string]interface{}, accountID int, userID int) ([]models.Account, error)
	ProcessAccount(data []string) error
	FindByAccountID(id uint) (*models.Account, error)
}
