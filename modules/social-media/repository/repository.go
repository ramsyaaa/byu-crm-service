package repository

import "byu-crm-service/models"

type SocialMediaRepository interface {
	GetSocialMediaBySubject(subject_type string, subject_id uint) ([]models.SocialMedia, error)
	DeleteBySubject(subject_type string, subject_id uint) error
	Insert(socialMedias []models.SocialMedia) error
}
