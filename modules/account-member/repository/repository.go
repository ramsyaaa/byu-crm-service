package repository

import "byu-crm-service/models"

type AccountMemberRepository interface {
	GetBySubject(subject_type string, subject_id uint) ([]models.AccountMember, error)
	DeleteBySubject(subject_type string, subject_id uint) error
	Insert(accountMember []models.AccountMember) error
}
