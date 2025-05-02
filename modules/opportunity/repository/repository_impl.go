package repository

import (
	"byu-crm-service/models"
	"errors"
	"strconv"
	"strings"
	"time"

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

func (r *opportunityRepository) CreateOpportunity(requestBody map[string]string, userID int) (*models.Opportunity, error) {
	opportunity := models.Opportunity{
		OpportunityName: func(s string) *string { return &s }(requestBody["opportunity_name"]),
		AccountID: func(s string) *uint {
			if s == "" {
				return nil
			}
			val, err := strconv.ParseUint(s, 10, 32)
			if err != nil {
				return nil
			}
			uval := uint(val)
			return &uval
		}(requestBody["account_id"]),
		ContactID: func(s string) *uint {
			if s == "" {
				return nil
			}
			val, err := strconv.ParseUint(s, 10, 32)
			if err != nil {
				return nil
			}
			uval := uint(val)
			return &uval
		}(requestBody["contact_id"]),
		Description: func(s string) *string { return &s }(requestBody["description"]),
		OpenDate: func(s string) *time.Time {
			if s == "" {
				return nil
			}
			val, err := time.Parse("2006-01-02", s)
			if err != nil {
				return nil
			}
			return &val
		}(requestBody["open_date"]),
		CloseDate: func(s string) *time.Time {
			if s == "" {
				return nil
			}
			val, err := time.Parse("2006-01-02", s)
			if err != nil {
				return nil
			}
			return &val
		}(requestBody["close_date"]),
		Amount: func(s string) *string { return &s }(requestBody["amount"]),
		CreatedBy: func(s string) *uint {
			if s == "" {
				return nil
			}
			val, err := strconv.ParseUint(s, 10, 32)
			if err != nil {
				return nil
			}
			uval := uint(val)
			return &uval
		}(requestBody["created_by"]),
	}

	if err := r.db.Create(&opportunity).Error; err != nil {
		return nil, err
	}

	var newOopportunity *models.Opportunity
	if err := r.db.Where("id = ?", opportunity.ID).First(&newOopportunity).Error; err != nil {
		return nil, err
	}

	return newOopportunity, nil
}

func (r *opportunityRepository) UpdateOpportunity(requestBody map[string]string, userID int, opportunityID int) (*models.Opportunity, error) {
	updateData := map[string]interface{}{}

	if val, ok := requestBody["opportunity_name"]; ok {
		updateData["opportunity_name"] = val
	}
	if val := requestBody["account_id"]; val != "" {
		if parsed, err := strconv.ParseUint(val, 10, 32); err == nil {
			updateData["account_id"] = uint(parsed)
		}
	}
	if val := requestBody["contact_id"]; val != "" {
		if parsed, err := strconv.ParseUint(val, 10, 32); err == nil {
			updateData["contact_id"] = uint(parsed)
		}
	}
	if val, ok := requestBody["description"]; ok {
		updateData["description"] = val
	}
	if val := requestBody["open_date"]; val != "" {
		if parsed, err := time.Parse("2006-01-02", val); err == nil {
			updateData["open_date"] = parsed
		}
	}
	if val := requestBody["close_date"]; val != "" {
		if parsed, err := time.Parse("2006-01-02", val); err == nil {
			updateData["close_date"] = parsed
		}
	}
	if val, ok := requestBody["amount"]; ok {
		updateData["amount"] = val
	}

	// Update data berdasarkan ID
	if err := r.db.Model(&models.Opportunity{}).Where("id = ?", opportunityID).Updates(updateData).Error; err != nil {
		return nil, err
	}

	// Ambil kembali data setelah update
	var updatedOpportunity *models.Opportunity
	if err := r.db.Where("id = ?", opportunityID).First(&updatedOpportunity).Error; err != nil {
		return nil, err
	}

	return updatedOpportunity, nil
}
