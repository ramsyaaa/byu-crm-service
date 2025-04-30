package repository

import (
	"byu-crm-service/models"
	"errors"
	"strings"

	"gorm.io/gorm"
)

type opportunityRepository struct {
	db *gorm.DB
}

func NewOpportunityRepository(db *gorm.DB) OpportunityRepository {
	return &opportunityRepository{db: db}
}

func (r *opportunityRepository) GetAllOpportunities(
	limit int,
	paginate bool,
	page int,
	filters map[string]string,
	userRole string,
	territoryID int,
) ([]models.Opportunity, int64, error) {
	var opportunity []models.Opportunity
	var total int64

	query := r.db.Model(&models.Opportunity{})

	// Apply search filter
	if search, exists := filters["search"]; exists && search != "" {
		searchTokens := strings.Fields(search)
		for _, token := range searchTokens {
			query = query.Where("opportunity.opportunity_name LIKE ?", "%"+token+"%")
		}
	}

	// Apply date range filter
	if startDate, exists := filters["start_date"]; exists && startDate != "" {
		query = query.Where("opportunity.created_at >= ?", startDate)
	}
	if endDate, exists := filters["end_date"]; exists && endDate != "" {
		query = query.Where("opportunity.created_at <= ?", endDate)
	}

	// Get total count before applying pagination
	query.Count(&total)

	// Apply ordering
	orderBy := filters["order_by"]
	order := filters["order"]
	if orderBy != "" && order != "" {
		query = query.Order(orderBy + " " + order)
	}

	// Apply pagination
	if paginate {
		offset := (page - 1) * limit
		query = query.Limit(limit).Offset(offset)
	} else if limit > 0 {
		query = query.Limit(limit)
	}

	// Preload relationships
	err := query.
		Preload("Account").
		Preload("Contact").
		Preload("User").
		Find(&opportunity).Error

	return opportunity, total, err
}

func (r *opportunityRepository) FindByOpportunityID(id uint, userRole string, territoryID uint, userID uint) (*models.Opportunity, error) {
	var opportunity models.Opportunity

	query := r.db.
		Model(&models.Opportunity{}).
		Preload("User").
		Preload("Contact").
		Preload("Account").
		Where("opportunities.id = ?", id)

	err := query.First(&opportunity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &opportunity, nil
}
