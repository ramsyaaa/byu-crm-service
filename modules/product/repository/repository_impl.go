package repository

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/product/response"
	"errors"
	"strconv"
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

func (r *productRepository) FindByProductID(id uint) (*models.Product, error) {
	var product models.Product

	query := r.db.
		Model(&models.Product{}).
		Preload("Eligibility", "subject_type = ? AND subject_id = ?", "App\\Models\\Product", id).
		Where("products.id = ?", id)

	err := query.First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) CreateProduct(requestBody map[string]string) (*models.Product, error) {
	product := models.Product{
		ProductName:     func(s string) *string { return &s }(requestBody["product_name"]),
		ProductType:     func(s string) *string { return &s }(requestBody["product_type"]),
		ProductCategory: func(s string) *string { return &s }(requestBody["product_category"]),
		Description:     func(s string) *string { return &s }(requestBody["description"]),
		Bid:             func(s string) *string { return &s }(requestBody["bid"]),
		Price: func(s string) *string {
			cleaned := strings.Map(func(r rune) rune {
				if r >= '0' && r <= '9' {
					return r
				}
				return -1
			}, s)

			return &cleaned
		}(requestBody["price"]),

		KeyVisual:      func(s string) *string { return &s }(requestBody["key_visual"]),
		AdditionalFile: func(s string) *string { return &s }(requestBody["additional_file"]),
		QuotaValue: func(s string) *float32 {
			if f, err := strconv.ParseFloat(s, 32); err == nil {
				val := float32(f)
				return &val
			}
			return nil
		}(requestBody["quota_value"]),
		ValidityValue: func(s string) *float32 {
			if f, err := strconv.ParseFloat(s, 32); err == nil {
				val := float32(f)
				return &val
			}
			return nil
		}(requestBody["validity_value"]),
		ValidityUnit: func(s string) *string { return &s }(requestBody["validity_unit"]),
	}

	if err := r.db.Create(&product).Error; err != nil {
		return nil, err
	}

	var newProduct *models.Product
	if err := r.db.Where("id = ?", product.ID).First(&newProduct).Error; err != nil {
		return nil, err
	}

	return newProduct, nil
}

func (r *productRepository) UpdateProduct(requestBody map[string]string, productID int) (*models.Product, error) {
	// Normalisasi harga
	cleanedPrice := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, requestBody["price"])

	product := models.Product{
		ProductName:     func(s string) *string { return &s }(requestBody["product_name"]),
		ProductType:     func(s string) *string { return &s }(requestBody["product_type"]),
		ProductCategory: func(s string) *string { return &s }(requestBody["product_category"]),
		Description:     func(s string) *string { return &s }(requestBody["description"]),
		Bid:             func(s string) *string { return &s }(requestBody["bid"]),
		Price:           &cleanedPrice,
		KeyVisual:       func(s string) *string { return &s }(requestBody["key_visual"]),
		AdditionalFile:  func(s string) *string { return &s }(requestBody["additional_file"]),
		QuotaValue: func(s string) *float32 {
			if f, err := strconv.ParseFloat(s, 32); err == nil {
				val := float32(f)
				return &val
			}
			return nil
		}(requestBody["quota_value"]),
		ValidityValue: func(s string) *float32 {
			if f, err := strconv.ParseFloat(s, 32); err == nil {
				val := float32(f)
				return &val
			}
			return nil
		}(requestBody["validity_value"]),
		ValidityUnit: func(s string) *string { return &s }(requestBody["validity_unit"]),
	}

	if err := r.db.Model(&models.Product{}).Where("id = ?", productID).Updates(&product).Error; err != nil {
		return nil, err
	}

	var updatedProduct models.Product
	if err := r.db.Where("id = ?", productID).First(&updatedProduct).Error; err != nil {
		return nil, err
	}

	return &updatedProduct, nil
}

func (r *productRepository) DeleteByAccountID(accountID uint) error {
	return r.db.Where("account_id = ?", accountID).Delete(&models.AccountProduct{}).Error
}

func (r *productRepository) Insert(productAccounts []models.AccountProduct) error {
	return r.db.Create(&productAccounts).Error
}
