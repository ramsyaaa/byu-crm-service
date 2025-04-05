package repository

import "byu-crm-service/models"

type AccountTypeSchoolDetailRepository interface {
	GetByAccountID(account_id uint) (*models.AccountTypeSchoolDetail, error)
	DeleteByAccountID(account_id uint) error
	Insert(accountTypeSchoolDetail *models.AccountTypeSchoolDetail) error
}
