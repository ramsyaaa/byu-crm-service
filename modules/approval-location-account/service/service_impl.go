package service

import (
	"byu-crm-service/models"
	accountRepo "byu-crm-service/modules/account/repository"
	"byu-crm-service/modules/approval-location-account/repository"
	"byu-crm-service/modules/approval-location-account/response"
	"fmt"
)

type approvalLocationAccountService struct {
	repo        repository.ApprovalLocationAccountRepository
	accountRepo accountRepo.AccountRepository
}

func NewApprovalLocationAccountService(repo repository.ApprovalLocationAccountRepository, accountRepo accountRepo.AccountRepository) ApprovalLocationAccountService {
	return &approvalLocationAccountService{repo: repo, accountRepo: accountRepo}
}

func (s *approvalLocationAccountService) GetAllApprovalRequest(limit int, paginate bool, page int, filters map[string]string, userRole string, territoryID int, userID int) ([]response.ApprovalLocationAccountResponse, int64, error) {
	return s.repo.GetAllApprovalRequest(limit, paginate, page, filters, userRole, territoryID, userID)
}

func (s *approvalLocationAccountService) CreateApprovalRequest(requestBody map[string]interface{}, userID int, accountID int) (*models.ApprovalLocationAccount, error) {
	// Use getStringValue to safely handle nil values and type conversions
	requestData := map[string]string{
		"account_id": getStringValue(accountID),
		"user_id":    getStringValue(userID),
		"longitude":  getStringValue(requestBody["longitude"]),
		"latitude":   getStringValue(requestBody["latitude"]),
		"status":     "0",
	}

	createApprovalRequest, err := s.repo.CreateApprovalRequest(requestData)
	if err != nil {
		return nil, err
	}

	return createApprovalRequest, nil
}

func (s *approvalLocationAccountService) FindByID(id uint) (*response.ApprovalLocationAccountResponse, error) {
	approvalRequest, err := s.repo.FindByID(uint(id))
	if err != nil {
		return nil, err
	}

	account, err := s.accountRepo.FindByAccountID(*approvalRequest.AccountID, "Super-Admin", 0, 0)

	if err != nil {
		return nil, err
	}

	if account.Latitude != nil && account.Longitude != nil {
		approvalRequest.LatestLatitude = account.Latitude
		approvalRequest.LatestLongitude = account.Longitude
	} else {
		approvalRequest.Latitude = nil
		approvalRequest.Longitude = nil
	}

	return approvalRequest, nil
}

func (s *approvalLocationAccountService) UpdateFields(id uint, fields map[string]interface{}) error {
	return s.repo.UpdateFields(id, fields)
}

func getStringValue(val interface{}) string {
	if val == nil {
		return ""
	}
	if str, ok := val.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", val)
}
