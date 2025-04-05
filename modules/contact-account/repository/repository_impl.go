package repository

import (
	"byu-crm-service/models"

	"gorm.io/gorm"
)

type contactAccountRepository struct {
	db *gorm.DB
}

func NewContactAccountRepository(db *gorm.DB) ContactAccountRepository {
	return &contactAccountRepository{db: db}
}

func (r *contactAccountRepository) GetByAccountID(account_id uint) ([]models.ContactAccount, error) {
	var contactAccounts []models.ContactAccount

	if err := r.db.Where("account_id = ?", account_id).Find(&contactAccounts).Error; err != nil {
		return nil, err
	}

	return contactAccounts, nil
}

func (r *contactAccountRepository) DeleteByAccountID(accountID uint) error {
	return r.db.Where("account_id = ?", accountID).Delete(&models.ContactAccount{}).Error
}

func (r *contactAccountRepository) Insert(contactAccounts []models.ContactAccount) error {
	return r.db.Create(&contactAccounts).Error
}
