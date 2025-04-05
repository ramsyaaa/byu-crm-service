package repository

import (
	"byu-crm-service/models"

	"gorm.io/gorm"
)

type socialMediaRepository struct {
	db *gorm.DB
}

func NewSocialMediaRepository(db *gorm.DB) SocialMediaRepository {
	return &socialMediaRepository{db: db}
}

func (r *socialMediaRepository) GetSocialMediaBySubject(subject_type string, subject_id uint) ([]models.SocialMedia, error) {
	var socialMedias []models.SocialMedia

	if err := r.db.Where("subject_type = ? AND subject_id = ?", subject_type, subject_id).Find(&socialMedias).Error; err != nil {
		return nil, err
	}

	return socialMedias, nil
}

func (r *socialMediaRepository) DeleteBySubject(subject_type string, subject_id uint) error {
	return r.db.Where("subject_type = ? AND subject_id = ?", subject_type, subject_id).
		Delete(&models.SocialMedia{}).Error
}

func (r *socialMediaRepository) Insert(socialMedias []models.SocialMedia) error {
	return r.db.Create(&socialMedias).Error
}
