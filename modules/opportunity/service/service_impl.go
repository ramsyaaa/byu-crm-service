package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/opportunity/repository"
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
