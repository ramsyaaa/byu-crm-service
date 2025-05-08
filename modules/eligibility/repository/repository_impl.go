package repository

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/eligibility/response"
	"encoding/json"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type eligibilityRepository struct {
	db *gorm.DB
}

func NewEligibilityRepository(db *gorm.DB) EligibilityRepository {
	return &eligibilityRepository{db: db}
}

func (s *eligibilityRepository) GetEligibilities(subjectType string, categories []string, types []string, locationFilter response.LocationFilter) ([]response.Eligibility, error) {
	var results []response.Eligibility

	query := s.db.Model(&response.Eligibility{}).Where("subject_type = ?", subjectType)

	// Filter categories dengan AND
	for _, cat := range categories {
		query = query.Where("JSON_CONTAINS(categories, ?)", "\""+cat+"\"")
	}

	// Filter types dengan AND
	for _, typ := range types {
		query = query.Where("JSON_CONTAINS(types, ?)", "\""+typ+"\"")
	}

	// Gabungkan semua kondisi lokasi dalam satu OR besar
	var locationConds []string
	var locationArgs []interface{}

	locationMap := map[string][]string{
		"areas":     locationFilter.Areas,
		"regions":   locationFilter.Regions,
		"branches":  locationFilter.Branches,
		"clusters":  locationFilter.Clusters,
		"cities":    locationFilter.Cities,
		"districts": locationFilter.Districts,
	}

	for key, values := range locationMap {
		for _, val := range values {
			escapedVal := strings.ReplaceAll(val, `"`, `\"`)
			jsonPath := fmt.Sprintf("$.%s", key)
			locationConds = append(locationConds, "JSON_CONTAINS(locations, ?, ?)")
			locationArgs = append(locationArgs, "\""+escapedVal+"\"", jsonPath)
		}
	}
	if len(locationConds) > 0 {
		// Satukan semua OR clause
		query = query.Where("("+strings.Join(locationConds, " OR ")+")", locationArgs...)
	}

	err := query.Find(&results).Error
	return results, err
}

func (s *eligibilityRepository) InsertEligibility(
	subjectType string,
	subjectID uint,
	categories []string,
	types []string,
	locations map[string][]string,
) error {
	// Hapus data lama
	if err := s.db.Where("subject_type = ? AND subject_id = ?", subjectType, subjectID).Delete(&models.Eligibility{}).Error; err != nil {
		return err
	}

	// Encode ke JSON string
	categoryJSON, err := json.Marshal(categories)
	if err != nil {
		return fmt.Errorf("failed to encode categories: %w", err)
	}

	typeJSON, err := json.Marshal(types)
	if err != nil {
		return fmt.Errorf("failed to encode types: %w", err)
	}

	locationJSON, err := json.Marshal(locations)
	if err != nil {
		return fmt.Errorf("failed to encode locations: %w", err)
	}

	// Simpan data baru
	newEligibility := models.Eligibility{
		SubjectType: subjectType,
		SubjectID:   subjectID,
		Categories:  string(categoryJSON),
		Types:       string(typeJSON),
		Locations:   string(locationJSON),
	}

	return s.db.Create(&newEligibility).Error
}
