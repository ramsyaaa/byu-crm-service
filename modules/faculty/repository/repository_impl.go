package repository

import (
	"byu-crm-service/models"
	"strings"

	"gorm.io/gorm"
)

type facultyRepository struct {
	db *gorm.DB
}

func NewFacultyRepository(db *gorm.DB) FacultyRepository {
	return &facultyRepository{db: db}
}

func (r *facultyRepository) GetAllFaculties(limit int, paginate bool, page int, filters map[string]string) ([]models.Faculty, int64, error) {
	var faculties []models.Faculty
	var total int64

	query := r.db.Model(&models.Faculty{})

	// Apply search filter
	if search, exists := filters["search"]; exists && search != "" {
		searchTokens := strings.Fields(search) // Tokenisasi input berdasarkan spasi
		for _, token := range searchTokens {
			query = query.Where(
				r.db.Where("faculties.name LIKE ?", "%"+token+"%"),
			)
		}
	}

	// Apply date range filter
	if startDate, exists := filters["start_date"]; exists && startDate != "" {
		query = query.Where("faculties.created_at >= ?", startDate)
	}
	if endDate, exists := filters["end_date"]; exists && endDate != "" {
		query = query.Where("faculties.created_at <= ?", endDate)
	}

	// Get total count before applying pagination
	query.Count(&total)

	// Apply ordering
	orderBy := filters["order_by"]
	order := filters["order"]
	query = query.Order(orderBy + " " + order)

	// Apply pagination
	if paginate {
		offset := (page - 1) * limit
		query = query.Limit(limit).Offset(offset)
	} else if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&faculties).Error
	return faculties, total, err
}

func (r *facultyRepository) GetFacultyByID(id int) (*models.Faculty, error) {
	var faculty models.Faculty
	err := r.db.First(&faculty, id).Error
	if err != nil {
		return nil, err
	}
	return &faculty, nil
}

func (r *facultyRepository) GetFacultyByName(name string) (*models.Faculty, error) {
	var faculty models.Faculty
	err := r.db.Where("name = ?", name).First(&faculty).Error
	if err != nil {
		return nil, err
	}
	return &faculty, nil
}

func (r *facultyRepository) CreateFaculty(faculty *models.Faculty) (*models.Faculty, error) {
	if err := r.db.Create(faculty).Error; err != nil {
		return nil, err
	}

	var createdFaculty models.Faculty
	if err := r.db.First(&createdFaculty, "id = ?", faculty.ID).Error; err != nil {
		return nil, err
	}

	return &createdFaculty, nil
}

func (r *facultyRepository) UpdateFaculty(faculty *models.Faculty, id int) (*models.Faculty, error) {
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
