package service

import "byu-crm-service/models"

type AccountService interface {
	GetAllAccounts(limit int, paginate bool, page int, filters map[string]string, userRole string, territoryID int, userID int, onlyUserPic bool, excludeVisited bool) ([]models.Account, int64, error)
	CreateAccount(requestBody map[string]interface{}, userID int) ([]models.Account, error)
	UpdateAccount(requestBody map[string]interface{}, accountID int, userRole string, territoryID int, userID int) ([]models.Account, error)
	ProcessAccount(data []string) error
	FindByAccountID(id uint, userRole string, territoryID uint, userID uint) (*models.Account, error)
	GetAccountVisitCounts(filters map[string]string, userRole string, territoryID int, userID int) (int64, int64, int64, error)
}
