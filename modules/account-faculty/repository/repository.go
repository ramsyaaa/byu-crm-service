package repository

import "byu-crm-service/models"

type AccountFacultyRepository interface {
	GetByAccountID(account_id uint) ([]models.AccountFaculty, error)
	DeleteByAccountID(account_id uint) error
	Insert(accountFaculty []models.AccountFaculty) error
}
