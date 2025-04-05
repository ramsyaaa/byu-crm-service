package service

import "byu-crm-service/models"

type AccountTypeSchoolDetailService interface {
	GetByAccountID(account_id uint) (*models.AccountTypeSchoolDetail, error)
	Insert(requestBody map[string]interface{}, account_id uint) (*models.AccountTypeSchoolDetail, error)
}
