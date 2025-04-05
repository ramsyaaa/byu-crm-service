package repository

import "byu-crm-service/models"

type AccountTypeCampusDetailRepository interface {
	GetByAccountID(account_id uint) (*models.AccountTypeCampusDetail, error)
	DeleteByAccountID(account_id uint) error
	Insert(accountTypeCampusDetail *models.AccountTypeCampusDetail) error
}
