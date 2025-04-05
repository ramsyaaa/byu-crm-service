package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/account-type-community-detail/repository"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

type accountTypeCommunityDetailService struct {
	repo repository.AccountTypeCommunityDetailRepository
}

func NewAccountTypeCommunityDetailService(repo repository.AccountTypeCommunityDetailRepository) AccountTypeCommunityDetailService {
	return &accountTypeCommunityDetailService{repo: repo}
}

func (s *accountTypeCommunityDetailService) GetByAccountID(account_id uint) (*models.AccountTypeCommunityDetail, error) {
	return s.repo.GetByAccountID(account_id)
}

func (s *accountTypeCommunityDetailService) Insert(requestBody map[string]interface{}, account_id uint) (*models.AccountTypeCommunityDetail, error) {
	// Delete existing social media for the given account_id
	if err := s.repo.DeleteByAccountID(account_id); err != nil {
		return nil, err
	}

	jsonGenderResult, err := arrayToJsonEncode(requestBody, []string{"gender", "percentage_gender"})
	if err != nil {
		jsonGenderResult = ""
	}
	jsonRangeAgeResult, err := arrayToJsonEncode(requestBody, []string{"age", "percentage_age"})
	if err != nil {
		jsonRangeAgeResult = ""
	}
	jsonEducationalBackgroundResult, err := arrayToJsonEncode(requestBody, []string{"educational_background", "percentage_educational_background"})
	if err != nil {
		jsonEducationalBackgroundResult = ""
	}
	jsonProfessionResult, err := arrayToJsonEncode(requestBody, []string{"profession", "percentage_profession"})
	if err != nil {
		jsonProfessionResult = ""
	}
	jsonIncomeResult, err := arrayToJsonEncode(requestBody, []string{"income", "percentage_income"})
	if err != nil {
		jsonIncomeResult = ""
	}

	var accountTypeCommunityDetail = models.AccountTypeCommunityDetail{
		AccountID:                   &account_id,
		AccountSubtype:              getStringPointer(requestBody, "account_subtype"),
		Group:                       getStringPointer(requestBody, "group"),
		GroupName:                   getStringPointer(requestBody, "group_name"),
		RangeAge:                    &jsonRangeAgeResult,
		Gender:                      &jsonGenderResult,
		EducationalBackground:       &jsonEducationalBackgroundResult,
		Profession:                  &jsonProfessionResult,
		Income:                      &jsonIncomeResult,
		ProductService:              getStringPointer(requestBody, "product_service"),
		PotentialCollaborationItems: getStringPointer(requestBody, "potential_collaboration_items"),
		PotentionalCollaboration:    getStringPointer(requestBody, "potentional_collaboration"),
	}

	// Insert the new detail accounts into the database
	if err := s.repo.Insert(&accountTypeCommunityDetail); err != nil {
		return nil, err
	}

	return &accountTypeCommunityDetail, nil
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
