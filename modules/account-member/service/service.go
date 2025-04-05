package service

import "byu-crm-service/models"

type AccountMemberService interface {
	GetBySubject(subject_type string, subject_id uint) ([]models.AccountMember, error)
	Insert(requestBody map[string]interface{}, subject_type string, subject_id uint, key1 string, key2 string) ([]models.AccountMember, error)
}
