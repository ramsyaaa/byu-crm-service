package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/contact-account/response"
)

type ContactAccountService interface {
	GetAllContacts(limit int, paginate bool, page int, filters map[string]string, userRole string, territoryID int, accountID int) ([]response.ContactResponse, int64, error)
	FindByContactID(id uint, userRole string, territoryID uint) (*models.Contact, error)
	CreateContact(requestBody map[string]interface{}) (*models.Contact, error)
	UpdateContact(requestBody map[string]interface{}, contactID int) (*models.Contact, error)
	GetContactAccountByAccountID(account_id uint) ([]models.ContactAccount, error)
	InsertContactAccount(requestBody map[string]interface{}, account_id uint) ([]models.ContactAccount, error)
	InsertContactAccountByContactID(requestBody map[string]interface{}, contact_id uint) ([]models.ContactAccount, error)
}
