package repository

import (
	"byu-crm-service/models"

	"gorm.io/gorm"
)

type accountTypeCampusDetailRepository struct {
	db *gorm.DB
}

func NewAccountTypeCampusDetailRepository(db *gorm.DB) AccountTypeCampusDetailRepository {
	return &accountTypeCampusDetailRepository{db: db}
}

func (r *accountTypeCampusDetailRepository) GetByAccountID(account_id uint) (*models.AccountTypeCampusDetail, error) {
	var accountTypeCampusDetail models.AccountTypeCampusDetail

	if err := r.db.Where("account_id = ?", account_id).First(&accountTypeCampusDetail).Error; err != nil {
		return nil, err
	}

	return &accountTypeCampusDetail, nil
}

func (r *accountTypeCampusDetailRepository) DeleteByAccountID(account_id uint) error {
	return r.db.Where("account_id = ?", account_id).
		Delete(&models.AccountTypeCampusDetail{}).Error
}

func (r *accountTypeCampusDetailRepository) Insert(accountTypeCampusDetail *models.AccountTypeCampusDetail) error {
	return r.db.Create(&accountTypeCampusDetail).Error
}
