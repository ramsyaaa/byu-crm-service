package repository

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/communication/response"
)

type CommunicationRepository interface {
	GetAllCommunications(limit int, paginate bool, page int, filters map[string]string, accountID int) ([]models.Communication, int64, error)
	FindByCommunicationID(id uint) (*models.Communication, error)
	CreateCommunication(requestBody map[string]string) (*models.Communication, error)
	UpdateAccount(requestBody map[string]string, accountID int, userID int) ([]models.Account, error)
	GetFilteredAccounts(limit, page int, search, userRole, territoryID string) ([]models.Account, int, error)
	Create(account *models.Account) error
	UpdateFields(id uint, fields map[string]interface{}) error
	FindByAccountID(id uint, userRole string, territoryID uint, userID uint) (*response.AccountResponse, error)
	GetAccountVisitCounts(filters map[string]string, userRole string, territoryID int, userID int) (int64, int64, int64, error)
}
