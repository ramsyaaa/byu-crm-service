package repository

import (
	"byu-crm-service/modules/program/response"
)

type ProgramRepository interface {
	GetAllPrograms(limit int, paginate bool, page int, filters map[string]string, subjectIDs []int) ([]response.ProgramResponse, int64, error)
	// FindByProductID(id uint) (*models.Product, error)
	// CreateProduct(requestBody map[string]string) (*models.Product, error)
	// UpdateProduct(requestBody map[string]string, productID int) (*models.Product, error)
	// Insert(productAccounts []models.AccountProduct) error
	// DeleteByAccountID(accountID uint) error
}
