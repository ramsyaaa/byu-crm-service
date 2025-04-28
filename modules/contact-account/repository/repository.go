package repository

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/contact-account/response"
)

type ContactAccountRepository interface {
	GetAllContacts(limit int, paginate bool, page int, filters map[string]string, userRole string, territoryID int, AccountID int) ([]response.ContactResponse, int64, error)
	GetByAccountID(account_id uint) ([]models.ContactAccount, error)
	DeleteByAccountID(accountID uint) error
	Insert(contactAccounts []models.ContactAccount) error
}
