package repository

import (
	"byu-crm-service/models"
	"strings"

	"gorm.io/gorm"
)

type constantDataRepository struct {
	db *gorm.DB
}

func NewConstantDataRepository(db *gorm.DB) ConstantDataRepository {
	return &constantDataRepository{db: db}
}

func (r *constantDataRepository) GetAllConstants(limit int, paginate bool, page int, filters map[string]string, type_constant string) ([]models.ConstantData, int64, error) {
	var constant_data []models.ConstantData
	var total int64

	query := r.db.Model(&models.ConstantData{})

	// Apply search filter
	if search, exists := filters["search"]; exists && search != "" {
		searchTokens := strings.Fields(search) // Tokenisasi input berdasarkan spasi
		for _, token := range searchTokens {
			query = query.Where("constant_data.name LIKE ?", "%"+token+"%")
		}
	}

	// Apply date range filter
	if startDate, exists := filters["start_date"]; exists && startDate != "" {
		query = query.Where("constant_data.created_at >= ?", startDate)
	}
	if endDate, exists := filters["end_date"]; exists && endDate != "" {
		query = query.Where("constant_data.created_at <= ?", endDate)
	}

	// Apply type_constant filter
	if type_constant != "" {
		query = query.Where("constant_data.type = ?", type_constant)
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
	err = query.Find(&constant_data).Error
	return constant_data, total, err
}

func (r *constantDataRepository) CreateConstant(constantData models.ConstantData) (models.ConstantData, error) {
	err := r.db.Create(&constantData).Error
	return constantData, err
}

func (r *constantDataRepository) GetConstantByTypeAndValue(type_constant string, value string) (models.ConstantData, error) {
	var constantData models.ConstantData
	err := r.db.Where("type = ? AND value = ?", type_constant, value).First(&constantData).Error
	if err != nil {
		return constantData, err
	}
	return constantData, nil
}
