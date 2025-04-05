package repository

import "byu-crm-service/models"

type AccountTypeCommunityDetailRepository interface {
	GetByAccountID(account_id uint) (*models.AccountTypeCommunityDetail, error)
	DeleteByAccountID(account_id uint) error
	Insert(accountTypeCommunityDetail *models.AccountTypeCommunityDetail) error
}
