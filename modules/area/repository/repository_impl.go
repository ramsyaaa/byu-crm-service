package repository

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/area/response"
	"strings"

	"gorm.io/gorm"
)

type areaRepository struct {
	db *gorm.DB
}

func NewAreaRepository(db *gorm.DB) AreaRepository {
	return &areaRepository{db: db}
}

func (r *areaRepository) GetAllAreas(limit int, paginate bool, page int, filters map[string]string) ([]response.AreaResponse, int64, error) {
	var areas []response.AreaResponse
	var total int64

	query := r.db.Model(&models.Area{})

	// Apply search filter
	if search, exists := filters["search"]; exists && search != "" {
		searchTokens := strings.Fields(search) // Tokenisasi input berdasarkan spasi
		for _, token := range searchTokens {
			query = query.Where(
				r.db.Where("areas.name LIKE ?", "%"+token+"%"),
			)
		}
	}

	// Apply date range filter
	if startDate, exists := filters["start_date"]; exists && startDate != "" {
		query = query.Where("areas.created_at >= ?", startDate)
	}
	if endDate, exists := filters["end_date"]; exists && endDate != "" {
		query = query.Where("areas.created_at <= ?", endDate)
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

	err := query.Find(&areas).Error
	return areas, total, err
}

func (r *areaRepository) GetAreaByID(id int) (*response.AreaResponse, error) {
	var area models.Area
	err := r.db.First(&area, id).Error
	if err != nil {
		return nil, err
	}

	areaResponse := &response.AreaResponse{
		ID:   area.ID,
		Name: area.Name,
	}

	return areaResponse, nil
}

func (r *areaRepository) GetAreaByName(name string) (*response.AreaResponse, error) {
	var area models.Area
	err := r.db.Where("name = ?", name).First(&area).Error
	if err != nil {
		return nil, err
	}

	areaResponse := &response.AreaResponse{
		ID:   area.ID,
		Name: area.Name,
	}

	return areaResponse, nil
}

func (r *areaRepository) CreateArea(area *models.Area) (*response.AreaResponse, error) {
	if err := r.db.Create(area).Error; err != nil {
		return nil, err
	}

	var createdArea models.Area
	if err := r.db.First(&createdArea, "id = ?", area.ID).Error; err != nil {
		return nil, err
	}

	areaResponse := &response.AreaResponse{
		ID:   createdArea.ID,
		Name: createdArea.Name,
	}

	return areaResponse, nil
}

func (r *areaRepository) UpdateArea(area *models.Area, id int) (*response.AreaResponse, error) {
	var existingArea models.Area

	if err := r.db.First(&existingArea, "id = ?", id).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&existingArea).Updates(area).Error; err != nil {
		return nil, err
	}

	if err := r.db.First(&existingArea, "id = ?", id).Error; err != nil {
		return nil, err
	}

	areaResponse := &response.AreaResponse{
		ID:   existingArea.ID,
		Name: existingArea.Name,
	}

	return areaResponse, nil
}
