package repository

import (
	"byu-crm-service/models"
	"errors"

	"gorm.io/gorm"
)

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) GetFilteredAccounts(limit, page int, search, userRole, territoryID string) ([]models.Account, int, error) {
	var accounts []models.Account
	var totalRecords int64

	query := r.db.Model(&models.Account{})

	// Apply filters based on search
	if search != "" {
		query = query.Where("account_name LIKE ?", "%"+search+"%")
	}

	// Apply role-based territorial filters
	switch userRole {
	case "Area":
		query = query.Where("area_id = ?", territoryID)
	case "Regional":
		query = query.Where("region_id = ?", territoryID)
	case "Branch":
		query = query.Where("branch_id = ?", territoryID)
	}

	// Count total records
	if err := query.Count(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	if err := query.Limit(limit).Offset(offset).Find(&accounts).Error; err != nil {
		return nil, 0, err
	}

	return accounts, int(totalRecords), nil
}

func (r *accountRepository) FindByAccountName(account_name string) (*models.Account, error) {
	var account models.Account
	if err := r.db.Where("account_name = ?", account_name).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not found is not an error
		}
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) FindByAccountCode(code string) (*models.Account, error) {
	var account models.Account
	if err := r.db.Where("account_code = ?", code).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not found is not an error
		}
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) Create(account *models.Account) error {
	return r.db.Create(account).Error
}
