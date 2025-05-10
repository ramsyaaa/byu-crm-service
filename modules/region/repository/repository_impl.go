package repository

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/region/response"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type regionRepository struct {
	db *gorm.DB
}

func NewRegionRepository(db *gorm.DB) RegionRepository {
	return &regionRepository{db: db}
}

func (r *regionRepository) GetAllRegions(filters map[string]string, userRole string, territoryID int) ([]response.RegionResponse, int64, error) {
	var regions []response.RegionResponse
	var total int64

	query := r.db.Model(&models.Region{})

	// Apply search filter
	if search, exists := filters["search"]; exists && search != "" {
		searchTokens := strings.Fields(search) // Tokenisasi input berdasarkan spasi
		for _, token := range searchTokens {
			query = query.Where(
				r.db.Where("regions.name LIKE ?", "%"+token+"%"),
			)
		}
	}

	if userRole != "" && territoryID != 0 {
		switch userRole {
		case "Area":
			query = query.Where("regions.area_id = ?", territoryID)

		case "Region":
			query = query.Where("regions.id = ?", territoryID)

		case "Branch", "YAE", "Organic", "Buddies", "DS":
			var regionID int
			err := r.db.Table("branches").Select("region_id").Where("id = ?", territoryID).Scan(&regionID).Error
			if err != nil {
				return nil, 0, err
			}
			query = query.Where("regions.id = ?", regionID)

		case "Cluster", "Admin-Tap":
			var regionID int
			err := r.db.Raw(`
				SELECT regions.id FROM regions
				JOIN branches ON branches.region_id = regions.id
				JOIN clusters ON clusters.branch_id = branches.id
				WHERE clusters.id = ? LIMIT 1
			`, territoryID).Scan(&regionID).Error
			if err != nil {
				return nil, 0, err
			}
			query = query.Where("regions.id = ?", regionID)
		}
	}

	// Get total count before applying pagination
	query.Count(&total)

	err := query.Find(&regions).Error
	return regions, total, err
}

func (r *regionRepository) GetRegionByID(id int) (*response.RegionResponse, error) {
	var region models.Region
	err := r.db.First(&region, id).Error
	if err != nil {
		return nil, err
	}

	regionResponse := &response.RegionResponse{
		ID:   region.ID,
		Name: region.Name,
		AreaID: func() int {
			if region.AreaID != nil {
				areaID, _ := strconv.Atoi(*region.AreaID)
				return areaID
			}
			return 0
		}(),
	}

	return regionResponse, nil
}

func (r *regionRepository) GetRegionByName(name string) (*response.RegionResponse, error) {
	var region models.Region
	err := r.db.Where("name = ?", name).First(&region).Error
	if err != nil {
		return nil, err
	}

	regionResponse := &response.RegionResponse{
		ID:   region.ID,
		Name: region.Name,
		AreaID: func() int {
			if region.AreaID != nil {
				areaID, _ := strconv.Atoi(*region.AreaID)
				return areaID
			}
			return 0
		}(),
	}

	return regionResponse, nil
}

func (r *regionRepository) CreateRegion(region *models.Region) (*response.RegionResponse, error) {
	if err := r.db.Create(region).Error; err != nil {
		return nil, err
	}

	var createdRegion models.Region
	if err := r.db.First(&createdRegion, "id = ?", region.ID).Error; err != nil {
		return nil, err
	}

	regionResponse := &response.RegionResponse{
		ID:   createdRegion.ID,
		Name: createdRegion.Name,
		AreaID: func() int {
			if createdRegion.AreaID != nil {
				areaID, _ := strconv.Atoi(*createdRegion.AreaID)
				return areaID
			}
			return 0
		}(),
	}

	return regionResponse, nil
}

func (r *regionRepository) UpdateRegion(region *models.Region, id int) (*response.RegionResponse, error) {
	var existingRegion models.Region

	if err := r.db.First(&existingRegion, "id = ?", id).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&existingRegion).Updates(region).Error; err != nil {
		return nil, err
	}

	if err := r.db.First(&existingRegion, "id = ?", id).Error; err != nil {
		return nil, err
	}

	regionResponse := &response.RegionResponse{
		ID:   existingRegion.ID,
		Name: existingRegion.Name,
		AreaID: func() int {
			if existingRegion.AreaID != nil {
				areaID, _ := strconv.Atoi(*existingRegion.AreaID)
				return areaID
			}
			return 0
		}(),
	}

	return regionResponse, nil
}
