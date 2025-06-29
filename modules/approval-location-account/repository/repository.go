package repository

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/approval-location-account/response"
)

type ApprovalLocationAccountRepository interface {
	UpdateFields(id uint, fields map[string]interface{}) error
	CreateApprovalRequest(requestBody map[string]string) (*models.ApprovalLocationAccount, error)
	GetAllApprovalRequest(limit int, paginate bool, page int, filters map[string]string, userRole string, territoryID int, userID int) ([]response.ApprovalLocationAccountResponse, int64, error)
	FindByID(id uint) (*response.ApprovalLocationAccountResponse, error)
}
