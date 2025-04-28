package repository

import (
	"byu-crm-service/models"
	"strings"

	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) GetAllCategories(limit int, paginate bool, page int, filters map[string]string, module string) ([]models.Category, int64, error) {
	var categories []models.Category
	var total int64

	query := r.db.Model(&models.Category{})

	// Apply search filter
	if search, exists := filters["search"]; exists && search != "" {
		searchTokens := strings.Fields(search) // Tokenisasi input berdasarkan spasi
		for _, token := range searchTokens {
			query = query.Where("categories.name LIKE ?", "%"+token+"%")
		}
	}

	// Apply date range filter
	if startDate, exists := filters["start_date"]; exists && startDate != "" {
		query = query.Where("categories.created_at >= ?", startDate)
	}
	if endDate, exists := filters["end_date"]; exists && endDate != "" {
		query = query.Where("categories.created_at <= ?", endDate)
	}

	// Apply module_type filter
	if module != "" {
		query = query.Where("categories.module_type = ?", module)
	}

	// Get total count before applying pagination
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Apply ordering safely
	orderBy := filters["order_by"]
	order := filters["order"]
	if orderBy != "" && order != "" {
		query = query.Order(orderBy + " " + order)
	}

	// Apply pagination
	if paginate {
		offset := (page - 1) * limit
		query = query.Limit(limit).Offset(offset)
	} else if limit > 0 {
		query = query.Limit(limit)
	}

	// Eksekusi query akhir
	err = query.Find(&categories).Error
	return categories, total, err
}

func (r *categoryRepository) GetCategoryByID(id int) (*models.Category, error) {
	var category models.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) GetCategoryByName(name string) (*models.Category, error) {
	var category models.Category
	err := r.db.Where("name = ?", name).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) GetCategoryByNameAndModuleType(name string, moduleType string) (*models.Category, error) {
	var category models.Category
	err := r.db.Where("name = ? AND module_type = ?", name, moduleType).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) CreateCategory(requestBody models.Category) (*models.Category, error) {
	err := r.db.Create(&requestBody).Error
	if err != nil {
		return nil, err
	}
	return &requestBody, nil
}

func (r *categoryRepository) UpdateCategory(id int, updateCategory models.Category) (models.Category, error) {
	err := r.db.Model(&models.Category{}).Where("id = ?", id).Updates(updateCategory).Error
	if err != nil {
		return models.Category{}, err
	}

	updatedCategory, err := r.GetCategoryByID(id)
	if err != nil {
		return models.Category{}, err
	}

	return *updatedCategory, nil
}
