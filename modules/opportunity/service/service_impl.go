package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/opportunity/repository"
	"fmt"
)

type opportunityService struct {
	repo repository.OpportunityRepository
}

func NewOpportunityService(repo repository.OpportunityRepository) OpportunityService {
	return &opportunityService{repo: repo}
}

func (s *opportunityService) GetAllOpportunities(limit int, paginate bool, page int, filters map[string]string, userRole string, territoryID int) ([]models.Opportunity, int64, error) {
	return s.repo.GetAllOpportunities(limit, paginate, page, filters, userRole, territoryID)
}

func (s *opportunityService) FindByOpportunityID(id uint, userRole string, territoryID uint, userID uint) (*models.Opportunity, error) {
	return s.repo.FindByOpportunityID(id, userRole, territoryID, userID)
}

func (s *opportunityService) CreateOpportunity(requestBody map[string]interface{}, userID int) (*models.Opportunity, error) {
	// Use getStringValue to safely handle nil values and type conversions
	opportunityData := map[string]string{
		"opportunity_name": getStringValue(requestBody["opportunity_name"]),
		"account_id":       getStringValue(requestBody["account_id"]),
		"contact_id":       getStringValue(requestBody["contact_id"]),
		"description":      getStringValue(requestBody["description"]),
		"open_date":        getStringValue(requestBody["open_date"]),
		"close_date":       getStringValue(requestBody["close_date"]),
		"amount":           getStringValue(requestBody["amount"]),
		"created_by":       getStringValue(userID),
	}

	opportunity, err := s.repo.CreateOpportunity(opportunityData, userID)
	if err != nil {
		return nil, err
	}

	return opportunity, nil
}

func (s *opportunityService) UpdateOpportunity(requestBody map[string]interface{}, userID int, opportunityID int) (*models.Opportunity, error) {
	// Use getStringValue to safely handle nil values and type conversions
	opportunityData := map[string]string{
		"opportunity_name": getStringValue(requestBody["opportunity_name"]),
		"account_id":       getStringValue(requestBody["account_id"]),
		"contact_id":       getStringValue(requestBody["contact_id"]),
		"description":      getStringValue(requestBody["description"]),
		"open_date":        getStringValue(requestBody["open_date"]),
		"close_date":       getStringValue(requestBody["close_date"]),
		"amount":           getStringValue(requestBody["amount"]),
		"created_by":       getStringValue(userID),
	}

	opportunity, err := s.repo.UpdateOpportunity(opportunityData, userID, opportunityID)
	if err != nil {
		return nil, err
	}

	return opportunity, nil
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
