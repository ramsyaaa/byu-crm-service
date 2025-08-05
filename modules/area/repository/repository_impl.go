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

func (r *areaRepository) GetAllAreas(filters map[string]string, userRole string, territoryID int, withGeo bool) ([]response.AreaResponse, int64, error) {
	var areas []response.AreaResponse
	var total int64

	query := r.db.Model(&models.Area{})

	if withGeo {
		query = query.Select("id, name, geojson")
	} else {
		query = query.Select("id, name")
	}

	// Apply search filter
	if search, exists := filters["search"]; exists && search != "" {
		searchTokens := strings.Fields(search)
		for _, token := range searchTokens {
			query = query.Where(
				r.db.Where("areas.name LIKE ?", "%"+token+"%"),
			)
		}
	}

	// UserRole & TerritoryID Filtering
	if userRole != "" && territoryID != 0 {
		switch userRole {
		case "Area":
			query = query.Where("areas.id = ?", territoryID)

		case "Region":
			// Cari area_id dari regions
			var areaID int
			err := r.db.Table("regions").Select("area_id").Where("id = ?", territoryID).Scan(&areaID).Error
			if err != nil {
				return nil, 0, err
			}
			query = query.Where("areas.id = ?", areaID)

		case "Branch", "YAE", "Organic", "Buddies", "DS":
			// Cari area_id dari branches -> regions
			var areaID int
			err := r.db.Raw(`
				SELECT areas.id FROM areas
				JOIN regions ON regions.area_id = areas.id
				JOIN branches ON branches.region_id = regions.id
				WHERE branches.id = ? LIMIT 1
			`, territoryID).Scan(&areaID).Error
			if err != nil {
				return nil, 0, err
			}
			query = query.Where("areas.id = ?", areaID)

		case "Cluster", "Admin-Tap":
			// Cari area_id dari clusters -> branches -> regions
			var areaID int
			err := r.db.Raw(`
				SELECT areas.id FROM areas
				JOIN regions ON regions.area_id = areas.id
				JOIN branches ON branches.region_id = regions.id
				JOIN clusters ON clusters.branch_id = branches.id
				WHERE clusters.id = ? LIMIT 1
			`, territoryID).Scan(&areaID).Error
			if err != nil {
				return nil, 0, err
			}
			query = query.Where("areas.id = ?", areaID)
		}
	}

	// Get total count before applying pagination
	query.Count(&total)

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
		ID:      area.ID,
		Name:    area.Name,
		Geojson: *area.Geojson,
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
