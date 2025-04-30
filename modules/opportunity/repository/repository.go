package repository

import (
	"byu-crm-service/models"
)

type OpportunityRepository interface {
	GetAllOpportunities(limit int, paginate bool, page int, filters map[string]string, userRole string, territoryID int) ([]models.Opportunity, int64, error)
	FindByOpportunityID(id uint, userRole string, territoryID uint, userID uint) (*models.Opportunity, error)
	// GetAreaByName(name string) (*response.AreaResponse, error)
	// CreateArea(area *models.Area) (*response.AreaResponse, error)
	// UpdateArea(area *models.Area, id int) (*response.AreaResponse, error)
}
