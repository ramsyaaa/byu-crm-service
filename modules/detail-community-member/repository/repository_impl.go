package repository

import (
	"byu-crm-service/models"
	"errors"

	"gorm.io/gorm"
)

type detailCommunityMemberRepository struct {
	db *gorm.DB
}

func NewDetailCommunityMemberRepository(db *gorm.DB) DetailCommunityMemberRepository {
	return &detailCommunityMemberRepository{db: db}
}

func (r *detailCommunityMemberRepository) FindByPhone(code string, accountID uint) (*models.DetailCommunityMember, error) {
	var detailCommunityMember models.DetailCommunityMember
	if err := r.db.
		Where("phone = ? AND account_id = ?", code, accountID).
		First(&detailCommunityMember).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not found is not an error
		}
		return nil, err
	}
	return &detailCommunityMember, nil
}

func (r *detailCommunityMemberRepository) Create(detailCommunityMember *models.DetailCommunityMember) error {
	return r.db.Create(detailCommunityMember).Error
}
