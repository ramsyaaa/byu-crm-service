package repository

import (
	"byu-crm-service/models"

	"gorm.io/gorm"
)

type accountTypeSchoolDetailRepository struct {
	db *gorm.DB
}

func NewAccountTypeSchoolDetailRepository(db *gorm.DB) AccountTypeSchoolDetailRepository {
	return &accountTypeSchoolDetailRepository{db: db}
}

func (r *accountTypeSchoolDetailRepository) GetByAccountID(account_id uint) (*models.AccountTypeSchoolDetail, error) {
	var accountTypeSchoolDetail models.AccountTypeSchoolDetail

	if err := r.db.Where("account_id = ?", account_id).First(&accountTypeSchoolDetail).Error; err != nil {
		return nil, err
	}

	return &accountTypeSchoolDetail, nil
}

func (r *accountTypeSchoolDetailRepository) DeleteByAccountID(account_id uint) error {
	return r.db.Where("account_id = ?", account_id).
		Delete(&models.AccountTypeSchoolDetail{}).Error
}

func (r *accountTypeSchoolDetailRepository) Insert(accountTypeSchoolDetail *models.AccountTypeSchoolDetail) error {
	return r.db.Create(&accountTypeSchoolDetail).Error
}
