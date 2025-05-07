package repository

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/product/response"
	"strings"

	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) GetAllProducts(limit int, paginate bool, page int, filters map[string]string, subjectIDs []int) ([]response.ProductResponse, int64, error) {
	var products []response.ProductResponse
	var total int64

	query := r.db.Model(&models.Product{})

	if len(subjectIDs) > 0 {
		query = query.Where("products.id IN ?", subjectIDs)
	}

	// Apply search filter
	if search, exists := filters["search"]; exists && search != "" {
		searchTokens := strings.Fields(search)
		for _, token := range searchTokens {
			query = query.Where(
				r.db.Where("products.product_name LIKE ?", "%"+token+"%"),
			)
		}
	}

	// Filter by date range
	if startDate, exists := filters["start_date"]; exists && startDate != "" {
		query = query.Where("products.created_at >= ?", startDate)
	}
	if endDate, exists := filters["end_date"]; exists && endDate != "" {
		query = query.Where("products.created_at <= ?", endDate)
	}

	// Count total before pagination
	query.Count(&total)

	// Apply ordering
	orderBy := filters["order_by"]
	order := filters["order"]
	query = query.Order(orderBy + " " + order)

	// Pagination
	if paginate {
		offset := (page - 1) * limit
		query = query.Limit(limit).Offset(offset)
	} else if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&products).Error
	return products, total, err
}
