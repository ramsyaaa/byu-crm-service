package repository

import (
	"byu-crm-service/models"
	"errors"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type communicationRepository struct {
	db *gorm.DB
}

func NewCommunicationRepository(db *gorm.DB) CommunicationRepository {
	return &communicationRepository{db: db}
}

func (r *communicationRepository) GetAllCommunications(
	limit int,
	paginate bool,
	page int,
	filters map[string]string,
	accountID int,
) ([]models.Communication, int64, error) {
	var communications []models.Communication
	var total int64

	query := r.db.Model(&models.Communication{}).Where("communications.account_id = ?", accountID).
		Preload("MainCommunication").
		Preload("NextCommunication").
		Preload("Account").
		Preload("Contact").
		Preload("Opportunity")

	// Apply search filter
	if search, exists := filters["search"]; exists && search != "" {
		searchTokens := strings.Fields(search)
		for _, token := range searchTokens {
			query = query.Where(
				r.db.Where("communications.note LIKE ?", "%"+token+"%"),
			)
		}
	}

	// Filter by date range
	if startDate, exists := filters["start_date"]; exists && startDate != "" {
		query = query.Where("communications.created_at >= ?", startDate)
	}
	if endDate, exists := filters["end_date"]; exists && endDate != "" {
		query = query.Where("communications.created_at <= ?", endDate)
	}

	// Count total sebelum pagination
	query.Count(&total)

	// Apply ordering
	orderBy := filters["order_by"]
	order := filters["order"]
	query = query.Order(orderBy + " " + order)

	// Pagination
	if paginate {
		offset := (page - 1) * limit
		query = query.Limit(limit).Offset(offset)
	} else if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&communications).Error
	return communications, total, err
}

func (r *communicationRepository) FindByCommunicationID(id uint) (*models.Communication, error) {
	var communication models.Communication

	query := r.db.
		Model(&models.Communication{}).
		Preload("MainCommunication").
		Preload("NextCommunication").
		Preload("Account").
		Preload("Contact").
		Preload("Opportunity").
		Where("communications.id = ?", id)

	err := query.First(&communication).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &communication, nil
}

func (r *communicationRepository) CreateCommunication(requestBody map[string]string) (*models.Communication, error) {
	communication := models.Communication{
		CommunicationType:   func(s string) *string { return &s }(requestBody["communication_type"]),
		Note:                func(s string) *string { return &s }(requestBody["note"]),
		Date:                func() *time.Time { now := time.Now(); return &now }(),
		StatusCommunication: func(s string) *string { return &s }(requestBody["status_communication"]),
		OpportunityID: func(s string) *uint {
			if s == "" {
				return nil
			}
			val, err := strconv.ParseUint(s, 10, 32)
			if err != nil {
				return nil
			}
			uval := uint(val)
			return &uval
		}(requestBody["opportunity_id"]),
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
		MainCommunicationID: func(s string) *uint {
			if s == "" {
				return nil
			}
			val, err := strconv.ParseUint(s, 10, 32)
			if err != nil {
				return nil
			}
			uval := uint(val)
			return &uval
		}(requestBody["main_communication_id"]),
	}

	if err := r.db.Create(&communication).Error; err != nil {
		return nil, err
	}

	var newCommunication, _ = r.FindByCommunicationID(communication.ID)

	return newCommunication, nil
}

func (r *communicationRepository) UpdateFields(id uint, fields map[string]interface{}) error {
	return r.db.Model(&models.Communication{}).Where("id = ?", id).Updates(fields).Error
}
