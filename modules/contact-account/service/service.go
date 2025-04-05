package service

import "byu-crm-service/models"

type ContactAccountService interface {
	GetContactAccountByAccountID(account_id uint) ([]models.ContactAccount, error)
	InsertContactAccount(requestBody map[string]interface{}, account_id uint) ([]models.ContactAccount, error)
}
