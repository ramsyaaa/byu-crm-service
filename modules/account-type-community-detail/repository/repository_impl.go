package repository

import (
	"byu-crm-service/models"

	"gorm.io/gorm"
)

type accountTypeCommunityDetailRepository struct {
	db *gorm.DB
}

func NewAccountTypeCommunityDetailRepository(db *gorm.DB) AccountTypeCommunityDetailRepository {
	return &accountTypeCommunityDetailRepository{db: db}
}

func (r *accountTypeCommunityDetailRepository) GetByAccountID(account_id uint) (*models.AccountTypeCommunityDetail, error) {
	var accountTypeCommunityDetail models.AccountTypeCommunityDetail

	if err := r.db.Where("account_id = ?", account_id).First(&accountTypeCommunityDetail).Error; err != nil {
		return nil, err
	}

	return &accountTypeCommunityDetail, nil
}

func (r *accountTypeCommunityDetailRepository) DeleteByAccountID(account_id uint) error {
	return r.db.Where("account_id = ?", account_id).
		Delete(&models.AccountTypeCommunityDetail{}).Error
}

func (r *accountTypeCommunityDetailRepository) Insert(accountTypeCommunityDetail *models.AccountTypeCommunityDetail) error {
	return r.db.Create(&accountTypeCommunityDetail).Error
}
