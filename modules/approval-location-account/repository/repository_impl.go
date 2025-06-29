package repository

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/approval-location-account/response"
	"errors"
	"strconv"

	"gorm.io/gorm"
)

type approvalLocationAccountRepository struct {
	db *gorm.DB
}

func NewApprovalLocationAccountRepository(db *gorm.DB) ApprovalLocationAccountRepository {
	return &approvalLocationAccountRepository{db: db}
}

func (r *approvalLocationAccountRepository) GetAllApprovalRequest(limit int, paginate bool, page int, filters map[string]string, userRole string, territoryID int, userID int) ([]response.ApprovalLocationAccountResponse, int64, error) {
	var approval_location_accounts []response.ApprovalLocationAccountResponse
	var total int64

	query := r.db.Model(&models.ApprovalLocationAccount{}).Where("approval_location_accounts.status = ?", 0)

	// Filter by date range
	if startDate, exists := filters["start_date"]; exists && startDate != "" {
		query = query.Where("approval_location_accounts.created_at >= ?", startDate)
	}
	if endDate, exists := filters["end_date"]; exists && endDate != "" {
		query = query.Where("approval_location_accounts.created_at <= ?", endDate)
	}

	// Count total before pagination
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

	err := query.Find(&approval_location_accounts).Error
	if err != nil {
		return approval_location_accounts, total, err
	}
	return approval_location_accounts, total, nil
}

func (r *approvalLocationAccountRepository) CreateApprovalRequest(requestBody map[string]string) (*models.ApprovalLocationAccount, error) {
	approvalRequest := models.ApprovalLocationAccount{
		Longitude: func(s string) *string { return &s }(requestBody["longitude"]),
		Latitude:  func(s string) *string { return &s }(requestBody["latitude"]),
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
		UserID: func(s string) *uint {
			if s == "" {
				return nil
			}
			val, err := strconv.ParseUint(s, 10, 32)
			if err != nil {
				return nil
			}
			uval := uint(val)
			return &uval
		}(requestBody["user_id"]),
		Status: func(s string) *uint {
			if s == "" {
				return nil
			}
			val, err := strconv.ParseUint(s, 10, 32)
			if err != nil {
				return nil
			}
			uval := uint(val)
			return &uval
		}(requestBody["status"]),
	}

	if err := r.db.Create(&approvalRequest).Error; err != nil {
		return nil, err
	}

	var approvalRequestNew *models.ApprovalLocationAccount
	if err := r.db.Where("id = ?", approvalRequest.ID).First(&approvalRequestNew).Error; err != nil {
		return nil, err
	}

	return approvalRequestNew, nil
}

func (r *approvalLocationAccountRepository) FindByID(id uint) (*response.ApprovalLocationAccountResponse, error) {
	var approvalLocationAccount response.ApprovalLocationAccountResponse

	query := r.db.
		Model(&models.ApprovalLocationAccount{}).
		Where("approval_location_accounts.id = ?", id).
		Where("approval_location_accounts.status = ?", 0)

	err := query.First(&approvalLocationAccount).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &approvalLocationAccount, nil
}

func (r *approvalLocationAccountRepository) UpdateFields(id uint, fields map[string]interface{}) error {
	return r.db.Model(&models.ApprovalLocationAccount{}).Where("id = ?", id).Updates(fields).Error
}
