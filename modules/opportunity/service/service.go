package service

import (
	"byu-crm-service/models"
)

type OpportunityService interface {
	GetAllOpportunities(limit int, paginate bool, page int, filters map[string]string, userRole string, territoryID int) ([]models.Opportunity, int64, error)
	FindByOpportunityID(id uint, userRole string, territoryID uint, userID uint) (*models.Opportunity, error)
}
