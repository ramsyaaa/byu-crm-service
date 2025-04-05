package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/account-type-campus-detail/repository"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type accountTypeCampusDetailService struct {
	repo repository.AccountTypeCampusDetailRepository
}

func NewAccountTypeCampusDetailService(repo repository.AccountTypeCampusDetailRepository) AccountTypeCampusDetailService {
	return &accountTypeCampusDetailService{repo: repo}
}

func (s *accountTypeCampusDetailService) GetByAccountID(account_id uint) (*models.AccountTypeCampusDetail, error) {
	return s.repo.GetByAccountID(account_id)
}

func (s *accountTypeCampusDetailService) Insert(requestBody map[string]interface{}, account_id uint) (*models.AccountTypeCampusDetail, error) {
	// Delete existing social media for the given account_id
	if err := s.repo.DeleteByAccountID(account_id); err != nil {
		return nil, err
	}

	jsonOriginResult, err := arrayToJsonEncode(requestBody, []string{"origin", "percentage_origin"})
	if err != nil {
		jsonOriginResult = ""
	}
	jsonRangeAgeResult, err := arrayToJsonEncode(requestBody, []string{"age", "percentage_age"})
	if err != nil {
		jsonRangeAgeResult = ""
	}
	jsonOrganizationNameResult, err := arrayToJsonEncode(requestBody, []string{"organization_name"})
	if err != nil {
		jsonOrganizationNameResult = ""
	}
	jsonPreferenceTechnologiesResult, err := encodeToJSON(requestBody, "preference_technologies")
	if err != nil {
		jsonPreferenceTechnologiesResult = ""
	}
	jsonMemberNeedsResult, err := encodeToJSON(requestBody, "member_needs")
	if err != nil {
		jsonMemberNeedsResult = ""
	}
	jsonItInfrastructuresResult, err := encodeToJSON(requestBody, "it_infrastructures")
	if err != nil {
		jsonItInfrastructuresResult = ""
	}
	jsonDigitalCollaborationsResult, err := encodeToJSON(requestBody, "digital_collaborations")
	if err != nil {
		jsonDigitalCollaborationsResult = ""
	}
	jsonProgramIdentificationResult, err := encodeToJSON(requestBody, "program_identification")
	if err != nil {
		jsonProgramIdentificationResult = ""
	}
	jsonUniversityRankResult, err := arrayToJsonEncode(requestBody, []string{"year_rank", "rank"})
	if err != nil {
		jsonUniversityRankResult = ""
	}
	jsonFocusProgramStudyResult, err := arrayToJsonEncode(requestBody, []string{"program_study"})
	if err != nil {
		jsonFocusProgramStudyResult = ""
	}

	var byodValue uint = 0 // Default 0 jika nil atau tidak valid
	if byod, exists := requestBody["byod"]; exists {
		if byodStr, ok := byod.(string); ok { // Pastikan `byod` bertipe string
			parsedByod, err := strconv.ParseUint(byodStr, 10, 64) // Konversi string ke uint64
			if err == nil {
				byodValue = uint(parsedByod) // Konversi ke uint
			}
		}
	}

	var accountTypeCampusDetail = models.AccountTypeCampusDetail{
		AccountID:                &account_id,
		RangeAge:                 &jsonRangeAgeResult,
		Origin:                   &jsonOriginResult,
		OrganizationName:         &jsonOrganizationNameResult,
		PreferenceTechnologies:   &jsonPreferenceTechnologiesResult,
		MemberNeeds:              &jsonMemberNeedsResult,
		ItInfrastructures:        &jsonItInfrastructuresResult,
		DigitalCollaborations:    &jsonDigitalCollaborationsResult,
		Byod:                     &byodValue,
		AccessTechnology:         getStringPointer(requestBody, "access_technology"),
		CampusAdministrationApp:  getStringPointer(requestBody, "campus_administration_app"),
		PotentionalCollaboration: getStringPointer(requestBody, "potentional_collaboration"),
		UniversityRank:           &jsonUniversityRankResult,
		FocusProgramStudy:        &jsonFocusProgramStudyResult,
		ProgramIdentification:    &jsonProgramIdentificationResult,
	}

	// Insert the new detail accounts into the database
	if err := s.repo.Insert(&accountTypeCampusDetail); err != nil {
		return nil, err
	}

	return &accountTypeCampusDetail, nil
}

func arrayToJsonEncode(requestBody map[string]interface{}, keys []string) (string, error) {
	// Pastikan keys tidak kosong
	if len(keys) == 0 {
		return "", errors.New("keys cannot be empty")
	}

	// Cek apakah key pertama tersedia dalam requestBody
	mainKey := keys[0]
	rawData, exists := requestBody[mainKey]
	if !exists {
		return "", fmt.Errorf("%s is missing", mainKey)
	}

	// Variabel untuk menyimpan hasil konversi
	var result []map[string]interface{}

	switch data := rawData.(type) {
	case []string:
		// Jika data berupa array string
		for i := range data {
			entry := make(map[string]interface{})
			for _, key := range keys {
				// Ambil value untuk setiap key
				if rawValues, ok := requestBody[key]; ok {
					// Cek jika nilai berupa []string
					if values, isArray := rawValues.([]string); isArray {
						if i < len(values) {
							entry[key] = values[i]
						} else {
							entry[key] = nil
						}
					} else {
						entry[key] = rawValues
					}
				}
			}
			result = append(result, entry)
		}
	case string:
		// Jika hanya satu string, buat satu entri
		entry := make(map[string]interface{})
		for _, key := range keys {
			if val, ok := requestBody[key]; ok {
				entry[key] = val
			} else {
				entry[key] = nil
			}
		}
		result = append(result, entry)
	default:
		return "", fmt.Errorf("invalid data type: %v", reflect.TypeOf(data))
	}

	// Encode ke JSON string
	jsonData, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func encodeToJSON(requestBody map[string]interface{}, key string) (string, error) {
	// Ambil data dari requestBody berdasarkan key
	data, exists := requestBody[key]
	if !exists {
		return "[]", nil // Jika tidak ada, kembalikan JSON kosong
	}

	// Encode ke JSON string
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func getStringPointer(data map[string]interface{}, key string) *string {
	if val, exists := data[key]; exists {
		if strVal, ok := val.(string); ok {
			return &strVal
		}
	}
	return nil // Jika tidak ada atau bukan string, kembalikan nil
}
