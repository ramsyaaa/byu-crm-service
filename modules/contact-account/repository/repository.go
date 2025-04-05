package repository

import "byu-crm-service/models"

type ContactAccountRepository interface {
	GetByAccountID(account_id uint) ([]models.ContactAccount, error)
	DeleteByAccountID(accountID uint) error
	Insert(contactAccounts []models.ContactAccount) error
}
