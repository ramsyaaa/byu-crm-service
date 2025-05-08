package repository

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/product/response"
)

type ProductRepository interface {
	GetAllProducts(limit int, paginate bool, page int, filters map[string]string, subjectIDs []int) ([]response.ProductResponse, int64, error)
	FindByProductID(id uint) (*models.Product, error)
	CreateProduct(requestBody map[string]string) (*models.Product, error)
	UpdateProduct(requestBody map[string]string, productID int) (*models.Product, error)
	Insert(productAccounts []models.AccountProduct) error
	DeleteByAccountID(accountID uint) error
}
