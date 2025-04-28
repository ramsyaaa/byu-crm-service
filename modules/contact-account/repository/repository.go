package repository

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/contact-account/response"
)

type ContactAccountRepository interface {
	GetAllContacts(limit int, paginate bool, page int, filters map[string]string, userRole string, territoryID int, AccountID int) ([]response.ContactResponse, int64, error)
	FindByContactID(id uint, userRole string, territoryID uint) (*models.Contact, error)
	CreateContact(requestBody map[string]string) (*models.Contact, error)
	UpdateContact(requestBody map[string]string, contactID int) (*models.Contact, error)
	GetByAccountID(account_id uint) ([]models.ContactAccount, error)
	DeleteByAccountID(accountID uint) error
	DeleteAccountByContactID(contactID uint) error
	Insert(contactAccounts []models.ContactAccount) error
}
