package repository

import (
	"byu-crm-service/modules/product/response"
)

type ProductRepository interface {
	GetAllProducts(limit int, paginate bool, page int, filters map[string]string, subjectIDs []int) ([]response.ProductResponse, int64, error)
}
