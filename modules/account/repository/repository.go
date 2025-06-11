package repository

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/account/response"
	"time"
)

type AccountRepository interface {
	FindByAccountCode(code string) (*models.Account, error)
	FindByAccountName(account_name string) (*models.Account, error)
	GetAllAccounts(limit int, paginate bool, page int, filters map[string]string, userRole string, territoryID int, userID int, onlyUserPic bool, excludeVisited bool) ([]response.AccountResponse, int64, error)
	CountAccount(userRole string, territoryID int, withGeoJson bool) (int64, map[string]int64, []map[string]interface{}, response.TerritoryInfo, error)
	CountByTerritories(userRole string, territoryID int, withGeoJson bool) ([]map[string]interface{}, error)
	CreateAccount(requestBody map[string]string, userID int) ([]models.Account, error)
	UpdateAccount(requestBody map[string]string, accountID int, userID int) ([]models.Account, error)
	Create(account *models.Account) error
	UpdateFields(id uint, fields map[string]interface{}) error
	FindByAccountID(id uint, userRole string, territoryID uint, userID uint) (*response.AccountResponse, error)
	GetAccountVisitCounts(filters map[string]string, userRole string, territoryID int, userID int) (int64, int64, int64, error)
	CheckAlreadyUpdateData(accountID int, userID int, clockInTime time.Time) (bool, error)
	CreateHistoryActivityAccount(userID, accountID uint, updateType string, subjectType *string, subjectID *uint) error
	UpdatePicMultipleAccounts(accountIDs []int, picID int) error
	UpdateAccountsPriority(accountIDs []int, priority string) error
	FindAccountsWithDifferentPic(accountIDs []int, userID int) ([]models.Account, error)
}
