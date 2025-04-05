package service

import "byu-crm-service/models"

type AccountTypeCommunityDetailService interface {
	GetByAccountID(account_id uint) (*models.AccountTypeCommunityDetail, error)
	Insert(requestBody map[string]interface{}, account_id uint) (*models.AccountTypeCommunityDetail, error)
}
