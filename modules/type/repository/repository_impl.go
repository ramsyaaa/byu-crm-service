package repository

import (
	"byu-crm-service/models"
	"strings"

	"gorm.io/gorm"
)

type typeRepository struct {
	db *gorm.DB
}

func NewTypeRepository(db *gorm.DB) TypeRepository {
	return &typeRepository{db: db}
}

func (r *typeRepository) GetAllTypes(limit int, paginate bool, page int, filters map[string]string, module string, category_name []string) ([]models.Type, int64, error) {
	var categories []models.Type
	var total int64

	query := r.db.Model(&models.Type{})

	var categoryIDs []uint
	if len(category_name) > 0 && module != "" {
		err := r.db.Model(&models.Category{}).
			Where("name IN ?", category_name).
			Where("module_type = ?", module).
			Pluck("id", &categoryIDs).Error
		if err != nil {
			return nil, 0, err
		}

		// Jika category_id ditemukan, filter types berdasarkan category_id
		if len(categoryIDs) > 0 {
			query = query.Where("category_id IN ?", categoryIDs)
		}
	}

	// Apply search filter
	if search, exists := filters["search"]; exists && search != "" {
		searchTokens := strings.Fields(search) // Tokenisasi input berdasarkan spasi
		for _, token := range searchTokens {
			query = query.Where("types.name LIKE ?", "%"+token+"%")
		}
	}

	// Apply date range filter
	if startDate, exists := filters["start_date"]; exists && startDate != "" {
		query = query.Where("types.created_at >= ?", startDate)
	}
	if endDate, exists := filters["end_date"]; exists && endDate != "" {
		query = query.Where("types.created_at <= ?", endDate)
	}

	// Apply module_type filter
	if module != "" {
		query = query.Where("types.module_type = ?", module)
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

func (r *typeRepository) GetTypeByID(id int) (models.Type, error) {
	var typeData models.Type
	err := r.db.Where("id = ?", id).First(&typeData).Error
	if err != nil {
		return models.Type{}, err
	}
	return typeData, nil
}

func (r *typeRepository) GetTypeByName(name string) (models.Type, error) {
	var typeData models.Type
	err := r.db.Where("name = ?", name).First(&typeData).Error
	if err != nil {
		return models.Type{}, err
	}
	return typeData, nil
}

func (r *typeRepository) GetTypeByNameAndModuleType(name string, moduleType string, categoryID int) (models.Type, error) {
	var category models.Type
	query := r.db.Where("name = ? AND module_type = ?", name, moduleType)

	if categoryID != 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	err := query.First(&category).Error
	if err != nil {
		return models.Type{}, err
	}
	return category, nil
}

func (r *typeRepository) CreateType(newType models.Type) (models.Type, error) {
	err := r.db.Create(&newType).Error
	if err != nil {
		return models.Type{}, err
	}
	return newType, nil
}

func (r *typeRepository) UpdateType(id int, updateType models.Type) (models.Type, error) {
	err := r.db.Model(&models.Type{}).Where("id = ?", id).Updates(updateType).Error
	if err != nil {
		return models.Type{}, err
	}

	return r.GetTypeByID(id)
}
