package service

import "byu-crm-service/models"

type SocialMediaService interface {
	GetSocialMediaBySubject(subject_type string, subject_id uint) ([]models.SocialMedia, error)
	InsertSocialMedia(requestBody map[string]interface{}, subject_type string, subject_id uint) ([]models.SocialMedia, error)
}
