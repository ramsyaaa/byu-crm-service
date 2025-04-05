package repository

import (
	"byu-crm-service/models"

	"gorm.io/gorm"
)

type accountMemberRepository struct {
	db *gorm.DB
}

func NewAccountMemberRepository(db *gorm.DB) AccountMemberRepository {
	return &accountMemberRepository{db: db}
}

func (r *accountMemberRepository) GetBySubject(subject_type string, subject_id uint) ([]models.AccountMember, error) {
	var accountMember []models.AccountMember

	if err := r.db.Where("subject_type = ? AND subject_id = ?", subject_type, subject_id).First(&accountMember).Error; err != nil {
		return nil, err
	}

	return accountMember, nil
}

func (r *accountMemberRepository) DeleteBySubject(subject_type string, subject_id uint) error {
	return r.db.Where("subject_type = ? AND subject_id = ?", subject_type, subject_id).
		Delete(&models.AccountMember{}).Error
}

func (r *accountMemberRepository) Insert(accountMember []models.AccountMember) error {
	return r.db.Create(&accountMember).Error
}
