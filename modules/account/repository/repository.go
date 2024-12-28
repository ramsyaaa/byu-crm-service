package repository

import "byu-crm-service/models"

type AccountRepository interface {
	FindByAccountCode(code string) (*models.Account, error)
	Create(account *models.Account) error
}
