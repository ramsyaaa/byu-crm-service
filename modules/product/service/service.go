package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/product/response"
)

type ProductService interface {
	GetAllProducts(limit int, paginate bool, page int, filters map[string]string, userRole string, territoryID int, userID int, accountID int) ([]response.ProductResponse, int64, error)
	FindByProductID(id uint) (*response.SingleProductResponse, error)
	CreateProduct(requestBody map[string]interface{}) (*models.Product, error)
	UpdateProduct(requestBody map[string]interface{}, productID int) (*models.Product, error)
	InsertProductAccount(requestBody map[string]interface{}, account_id uint) ([]models.AccountProduct, error)
}
