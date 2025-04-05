package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/social-media/repository"
	"errors"
	"fmt"
	"reflect"
)

type socialMediaService struct {
	repo repository.SocialMediaRepository
}

func NewSocialMediaService(repo repository.SocialMediaRepository) SocialMediaService {
	return &socialMediaService{repo: repo}
}

func (s *socialMediaService) GetSocialMediaBySubject(subject_type string, subject_id uint) ([]models.SocialMedia, error) {
	return s.repo.GetSocialMediaBySubject(subject_type, subject_id)
}

func (s *socialMediaService) InsertSocialMedia(requestBody map[string]interface{}, subject_type string, subject_id uint) ([]models.SocialMedia, error) {
	// Delete existing social media for the given subject_type and subject_id
	if err := s.repo.DeleteBySubject(subject_type, subject_id); err != nil {
		return nil, err
	}

	categories, exists := requestBody["category"]
	if !exists {
		return nil, errors.New("category is missing")
	}

	urls, exists := requestBody["url"]
	if !exists {
		return nil, errors.New("url is missing")
	}

	var DataCategory []string

	switch v := categories.(type) {
	case string: // Jika hanya satu nilai, ubah ke array
		DataCategory = append(DataCategory, v)
	case []string: // Jika sudah array string, langsung tambahkan ke DataCategory
		DataCategory = append(DataCategory, v...)
	case []interface{}: // Jika array tapi bertipe []interface{}
		for _, val := range v {
			strVal, ok := val.(string)
			if !ok {
				return nil, errors.New("contact_id contains non-string value")
			}
			DataCategory = append(DataCategory, strVal)
		}
	default:
		return nil, fmt.Errorf("invalid contact_id type: %v", reflect.TypeOf(categories))
	}

	var DataUrl []string

	switch v := urls.(type) {
	case string: // Jika hanya satu nilai, ubah ke array
		DataUrl = append(DataUrl, v)
	case []string: // Jika sudah array string, langsung tambahkan ke DataUrl
		DataUrl = append(DataUrl, v...)
	case []interface{}: // Jika array tapi bertipe []interface{}
		for _, val := range v {
			strVal, ok := val.(string)
			if !ok {
				return nil, errors.New("contact_id contains non-string value")
			}
			DataUrl = append(DataUrl, strVal)
		}
	default:
		return nil, fmt.Errorf("invalid contact_id type: %v", reflect.TypeOf(urls))
	}

	if len(DataCategory) != len(DataUrl) {
		return nil, errors.New("category and url length mismatch")
	}

	var socialMedias []models.SocialMedia
	for i := range DataCategory {
		socialMedias = append(socialMedias, models.SocialMedia{
			SubjectType: &subject_type,
			SubjectID:   &subject_id,
			Category:    &DataCategory[i],
			Url:         &DataUrl[i],
		})
	}

	// Insert the new contact accounts into the database
	if err := s.repo.Insert(socialMedias); err != nil {
		return nil, err
	}

	return socialMedias, nil
}
