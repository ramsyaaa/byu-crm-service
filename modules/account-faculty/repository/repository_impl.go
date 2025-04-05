package repository

import (
	"byu-crm-service/models"

	"gorm.io/gorm"
)

type accountFacultyRepository struct {
	db *gorm.DB
}

func NewAccountFacultyRepository(db *gorm.DB) AccountFacultyRepository {
	return &accountFacultyRepository{db: db}
}

func (r *accountFacultyRepository) GetByAccountID(account_id uint) ([]models.AccountFaculty, error) {
	var accountFaculty []models.AccountFaculty

	if err := r.db.Where("account_id = ?", account_id).First(&accountFaculty).Error; err != nil {
		return nil, err
	}

	return accountFaculty, nil
}

func (r *accountFacultyRepository) DeleteByAccountID(account_id uint) error {
	return r.db.Where("account_id = ?", account_id).
		Delete(&models.AccountFaculty{}).Error
}

func (r *accountFacultyRepository) Insert(accountFaculty []models.AccountFaculty) error {
	return r.db.Create(&accountFaculty).Error
}
