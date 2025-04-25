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

func (r *categoryRepository) GetFacultyByID(id int) (*models.Faculty, error) {
	var faculty models.Faculty
	err := r.db.First(&faculty, id).Error
	if err != nil {
		return nil, err
	}
	return &faculty, nil
}

func (r *categoryRepository) GetFacultyByName(name string) (*models.Faculty, error) {
	var faculty models.Faculty
	err := r.db.Where("name = ?", name).First(&faculty).Error
	if err != nil {
		return nil, err
	}
	return &faculty, nil
}

func (r *categoryRepository) CreateFaculty(faculty *models.Faculty) (*models.Faculty, error) {
	if err := r.db.Create(faculty).Error; err != nil {
		return nil, err
	}

	var createdFaculty models.Faculty
	if err := r.db.First(&createdFaculty, "id = ?", faculty.ID).Error; err != nil {
		return nil, err
	}

	return &createdFaculty, nil
}

func (r *categoryRepository) UpdateFaculty(faculty *models.Faculty, id int) (*models.Faculty, error) {
	var existingFaculty models.Faculty
	if err := r.db.First(&existingFaculty, "id = ?", id).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&existingFaculty).Updates(faculty).Error; err != nil {
		return nil, err
	}

	if err := r.db.First(&existingFaculty, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &existingFaculty, nil
}
