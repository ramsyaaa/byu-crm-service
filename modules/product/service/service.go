package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/product/response"
)

type ProductService interface {
	GetAllProducts(limit int, paginate bool, page int, filters map[string]string, userRole string, territoryID int, userID int, accountID int) ([]response.ProductResponse, int64, error)
	InsertProductAccount(requestBody map[string]interface{}, account_id uint) ([]models.AccountProduct, error)
}
