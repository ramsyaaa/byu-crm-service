package service

import "byu-crm-service/models"

type AccountTypeCampusDetailService interface {
	GetByAccountID(account_id uint) (*models.AccountTypeCampusDetail, error)
	Insert(requestBody map[string]interface{}, account_id uint) (*models.AccountTypeCampusDetail, error)
}
