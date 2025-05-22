package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/account/response"
	"time"
)

type AccountService interface {
	GetAllAccounts(limit int, paginate bool, page int, filters map[string]string, userRole string, territoryID int, userID int, onlyUserPic bool, excludeVisited bool) ([]response.AccountResponse, int64, error)
	CountAccount(userRole string, territoryID int) (int64, map[string]int64, []map[string]interface{}, response.TerritoryInfo, error)
	CreateAccount(requestBody map[string]interface{}, userID int) ([]models.Account, error)
	UpdateAccount(requestBody map[string]interface{}, accountID int, userRole string, territoryID int, userID int) ([]models.Account, error)
	UpdatePic(accountID int, userRole string, territoryID int, userID int) (*response.AccountResponse, error)
	UpdateFields(id uint, fields map[string]interface{}) error
	ProcessAccount(data []string) error
	CheckAlreadyUpdateData(accountID int, clockIn time.Time, userID int) (bool, error)
	CreateHistoryActivityAccount(userID, accountID uint, updateType string, subjectType *string, subjectID *uint) error
	FindByAccountID(id uint, userRole string, territoryID uint, userID uint) (*response.SingleAccountResponse, error)
	GetAccountVisitCounts(filters map[string]string, userRole string, territoryID int, userID int) (int64, int64, int64, error)
}
