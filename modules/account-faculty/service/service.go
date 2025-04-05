package service

import "byu-crm-service/models"

type AccountFacultyService interface {
	GetByAccountID(account_id uint) ([]models.AccountFaculty, error)
	Insert(requestBody map[string]interface{}, account_id uint) ([]models.AccountFaculty, error)
}
