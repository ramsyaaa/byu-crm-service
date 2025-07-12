package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/approval-location-account/response"
)

type ApprovalLocationAccountService interface {
	FindByUserIDAndAccountID(userID uint, accountID int, status int) (*models.ApprovalLocationAccount, error)
	UpdateFields(id uint, fields map[string]interface{}) error
	GetAllApprovalRequest(limit int, paginate bool, page int, filters map[string]string, userRole string, territoryID int, userIDs []int) ([]response.ApprovalLocationAccountResponse, int64, error)
	FindByID(id uint) (*response.ApprovalLocationAccountResponse, error)
	CreateApprovalRequest(requestBody map[string]interface{}, userID int, accountID int) (*models.ApprovalLocationAccount, error)
}
