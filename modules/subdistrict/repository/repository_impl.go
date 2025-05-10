package repository

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/subdistrict/response"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type subdistrictRepository struct {
	db *gorm.DB
}

func NewSubdistrictRepository(db *gorm.DB) SubdistrictRepository {
	return &subdistrictRepository{db: db}
}

func (r *subdistrictRepository) GetAllSubdistricts(filters map[string]string, userRole string, territoryID int) ([]response.SubdistrictResponse, int64, error) {
	var subdistricts []response.SubdistrictResponse
	var total int64

	query := r.db.Model(&models.Subdistrict{}).
		Joins("JOIN cities ON cities.id = subdistricts.city_id").
		Joins("JOIN clusters ON clusters.id = cities.cluster_id").
		Joins("JOIN branches ON branches.id = clusters.branch_id").
		Joins("JOIN regions ON regions.id = branches.region_id")

	// Filter berdasarkan role dan territory
	if userRole != "" && territoryID != 0 {
		switch userRole {
		case "Cluster", "Admin-Tap":
			query = query.Where("cities.cluster_id = ?", territoryID)

		case "Branch", "YAE", "Organic", "Buddies", "DS":
			query = query.Where("clusters.branch_id = ?", territoryID)

		case "Region":
			query = query.Where("branches.region_id = ?", territoryID)

		case "Area":
			query = query.Where("regions.area_id = ?", territoryID)
		}
	}

	// Apply search filter
	if search, exists := filters["search"]; exists && search != "" {
		searchTokens := strings.Fields(search)
		for _, token := range searchTokens {
			query = query.Where("subdistricts.name LIKE ?", "%"+token+"%")
		}
	}

	// Get total count before pagination
	query.Count(&total)

	err := query.Find(&subdistricts).Error
	return subdistricts, total, err
}

func (r *subdistrictRepository) GetSubdistrictByID(id int) (*response.SubdistrictResponse, error) {
	var subdistrict models.Subdistrict
	err := r.db.First(&subdistrict, id).Error
	if err != nil {
		return nil, err
	}

	subdistrictResponse := &response.SubdistrictResponse{
		ID:   subdistrict.ID,
		Name: subdistrict.Name,
		CityID: func() int {
			if subdistrict.CityID != nil {
				CityID, _ := strconv.Atoi(*subdistrict.CityID)
				return CityID
			}
			return 0
		}(),
	}

	return subdistrictResponse, nil
}

func (r *subdistrictRepository) GetSubdistrictByName(name string) (*response.SubdistrictResponse, error) {
	var subdistrict models.Subdistrict
	err := r.db.Where("name = ?", name).First(&subdistrict).Error
	if err != nil {
		return nil, err
	}

	subdistrictResponse := &response.SubdistrictResponse{
		ID:   subdistrict.ID,
		Name: subdistrict.Name,
		CityID: func() int {
			if subdistrict.CityID != nil {
				cityID, _ := strconv.Atoi(*subdistrict.CityID)
				return cityID
			}
			return 0
		}(),
	}

	return subdistrictResponse, nil
}

func (r *subdistrictRepository) CreateSubdistrict(subdistrict *models.Subdistrict) (*response.SubdistrictResponse, error) {
	if err := r.db.Create(subdistrict).Error; err != nil {
		return nil, err
	}

	var createdSubdistrict models.Subdistrict
	if err := r.db.First(&createdSubdistrict, "id = ?", subdistrict.ID).Error; err != nil {
		return nil, err
	}

	subdistrictResponse := &response.SubdistrictResponse{
		ID:   createdSubdistrict.ID,
		Name: createdSubdistrict.Name,
		CityID: func() int {
			if createdSubdistrict.CityID != nil {
				cityID, _ := strconv.Atoi(*createdSubdistrict.CityID)
				return cityID
			}
			return 0
		}(),
	}

	return subdistrictResponse, nil
}

func (r *subdistrictRepository) UpdateSubdistrict(subdistrict *models.Subdistrict, id int) (*response.SubdistrictResponse, error) {
	var existingSubdistrict models.Subdistrict

	if err := r.db.First(&existingSubdistrict, "id = ?", id).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&existingSubdistrict).Updates(subdistrict).Error; err != nil {
		return nil, err
	}

	if err := r.db.First(&existingSubdistrict, "id = ?", id).Error; err != nil {
		return nil, err
	}

	subdistrictResponse := &response.SubdistrictResponse{
		ID:   existingSubdistrict.ID,
		Name: existingSubdistrict.Name,
		CityID: func() int {
			if existingSubdistrict.CityID != nil {
				cityID, _ := strconv.Atoi(*existingSubdistrict.CityID)
				return cityID
			}
			return 0
		}(),
	}

	return subdistrictResponse, nil
}
