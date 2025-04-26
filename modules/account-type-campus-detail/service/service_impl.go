package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/account-type-campus-detail/repository"
	"encoding/json"
	"fmt"
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

func (s *accountTypeCampusDetailService) Insert(requestBody map[string]interface{}, accountID uint) (*models.AccountTypeCampusDetail, error) {
	// Delete existing detail by account_id
	if err := s.repo.DeleteByAccountID(accountID); err != nil {
		return nil, err
	}

	// Helper encode
	jsonOriginResult, _ := arrayToJsonEncode(requestBody, []string{"origin", "percentage_origin"})
	jsonRangeAgeResult, _ := arrayToJsonEncode(requestBody, []string{"age", "percentage_age"})
	jsonOrganizationNameResult, _ := arrayToJsonEncode(requestBody, []string{"organization_name"})
	jsonPreferenceTechnologiesResult, _ := encodeToJSON(requestBody, "preference_technologies")
	jsonMemberNeedsResult, _ := encodeToJSON(requestBody, "member_needs")
	jsonItInfrastructuresResult, _ := encodeToJSON(requestBody, "it_infrastructures")
	jsonDigitalCollaborationsResult, _ := encodeToJSON(requestBody, "digital_collaborations")
	jsonProgramIdentificationResult, _ := encodeToJSON(requestBody, "program_identification")
	jsonUniversityRankResult, _ := arrayToJsonEncode(requestBody, []string{"year_rank", "rank"})
	jsonFocusProgramStudyResult, _ := arrayToJsonEncode(requestBody, []string{"program_study"})

	byodValue := getUintFromString(requestBody, "byod")

	accountTypeCampusDetail := models.AccountTypeCampusDetail{
		AccountID:                &accountID,
		RangeAge:                 nullableString(jsonRangeAgeResult),
		Origin:                   nullableString(jsonOriginResult),
		OrganizationName:         nullableString(jsonOrganizationNameResult),
		PreferenceTechnologies:   nullableString(jsonPreferenceTechnologiesResult),
		MemberNeeds:              nullableString(jsonMemberNeedsResult),
		ItInfrastructures:        nullableString(jsonItInfrastructuresResult),
		DigitalCollaborations:    nullableString(jsonDigitalCollaborationsResult),
		Byod:                     &byodValue,
		AccessTechnology:         getStringPointer(requestBody, "access_technology"),
		CampusAdministrationApp:  getStringPointer(requestBody, "campus_administration_app"),
		PotentionalCollaboration: getStringPointer(requestBody, "potentional_collaboration"),
		UniversityRank:           nullableString(jsonUniversityRankResult),
		FocusProgramStudy:        nullableString(jsonFocusProgramStudyResult),
		ProgramIdentification:    nullableString(jsonProgramIdentificationResult),
	}

	if err := s.repo.Insert(&accountTypeCampusDetail); err != nil {
		return nil, err
	}

	return &accountTypeCampusDetail, nil
}

func arrayToJsonEncode(requestBody map[string]interface{}, keys []string) (string, error) {
	if len(keys) == 0 {
		return "", nil
	}

	// Ambil semua array
	dataArrays := make([][]interface{}, len(keys))
	for i, key := range keys {
		if val, ok := requestBody[key]; ok {
			if arr, ok := val.([]interface{}); ok {
				dataArrays[i] = arr
			} else {
				return "", fmt.Errorf("field %s harus berupa array", key)
			}
		} else {
			// Kalau field tidak ada, berarti kosong
			dataArrays[i] = []interface{}{}
		}
	}

	// Pastikan semua array panjangnya sama
	length := len(dataArrays[0])
	for _, arr := range dataArrays {
		if len(arr) != length {
			return "", fmt.Errorf("panjang array %v tidak konsisten", keys)
		}
	}

	// Gabungkan array jadi slice of map
	result := make([]map[string]interface{}, length)
	for i := 0; i < length; i++ {
		item := make(map[string]interface{})
		for j, key := range keys {
			item[key] = dataArrays[j][i]
		}
		result[i] = item
	}

	if len(result) == 0 {
		return "", nil
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func encodeToJSON(requestBody map[string]interface{}, key string) (string, error) {
	if val, ok := requestBody[key]; ok {
		jsonBytes, err := json.Marshal(val)
		if err != nil {
			return "", err
		}
		return string(jsonBytes), nil
	}
	return "", nil
}

func getStringPointer(requestBody map[string]interface{}, key string) *string {
	if val, ok := requestBody[key]; ok {
		if strVal, ok := val.(string); ok && strVal != "" {
			return &strVal
		}
	}
	return nil
}

func getUintFromString(requestBody map[string]interface{}, key string) uint {
	if val, ok := requestBody[key]; ok {
		if strVal, ok := val.(string); ok {
			if parsed, err := strconv.ParseUint(strVal, 10, 64); err == nil {
				return uint(parsed)
			}
		}
	}
	return 0
}

func nullableString(str string) *string {
	if str == "" {
		return nil
	}
	return &str
}
