package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/account-type-community-detail/repository"
	"encoding/json"
	"fmt"
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

func (s *accountTypeCommunityDetailService) Insert(requestBody map[string]interface{}, accountID uint) (*models.AccountTypeCommunityDetail, error) {
	// Delete existing detail by account_id
	if err := s.repo.DeleteByAccountID(accountID); err != nil {
		return nil, err
	}

	// Helper encode
	jsonGenderResult, _ := arrayToJsonEncode(requestBody, []string{"gender", "percentage_gender"})
	jsonRangeAgeResult, _ := arrayToJsonEncode(requestBody, []string{"age", "percentage_age"})
	jsonEducationalBackgroundResult, _ := arrayToJsonEncode(requestBody, []string{"educational_background", "percentage_educational_background"})
	jsonProfessionResult, _ := arrayToJsonEncode(requestBody, []string{"profession", "percentage_profession"})
	jsonIncomeResult, _ := arrayToJsonEncode(requestBody, []string{"income", "percentage_income"})

	accountTypeCommunityDetail := models.AccountTypeCommunityDetail{
		AccountID:                   &accountID,
		AccountSubtype:              getStringPointer(requestBody, "account_subtype"),
		Group:                       getStringPointer(requestBody, "group"),
		GroupName:                   getStringPointer(requestBody, "group_name"),
		RangeAge:                    nullableString(jsonRangeAgeResult),
		Gender:                      nullableString(jsonGenderResult),
		EducationalBackground:       nullableString(jsonEducationalBackgroundResult),
		Profession:                  nullableString(jsonProfessionResult),
		Income:                      nullableString(jsonIncomeResult),
		ProductService:              getStringPointer(requestBody, "product_service"),
		PotentialCollaborationItems: getStringPointer(requestBody, "potential_collaboration_items"),
		PotentionalCollaboration:    getStringPointer(requestBody, "potentional_collaboration"),
	}

	if err := s.repo.Insert(&accountTypeCommunityDetail); err != nil {
		return nil, err
	}

	return &accountTypeCommunityDetail, nil
}

func arrayToJsonEncode(requestBody map[string]interface{}, keys []string) (string, error) {
	if len(keys) == 0 {
		return "", nil
	}

	dataArrays := make([][]interface{}, len(keys))
	for i, key := range keys {
		if val, ok := requestBody[key]; ok {
			if arr, ok := val.([]interface{}); ok {
				dataArrays[i] = arr
			} else {
				return "", fmt.Errorf("field %s harus berupa array", key)
			}
		} else {
			dataArrays[i] = []interface{}{}
		}
	}

	length := len(dataArrays[0])
	for _, arr := range dataArrays {
		if len(arr) != length {
			return "", fmt.Errorf("panjang array %v tidak konsisten", keys)
		}
	}

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

func getStringPointer(requestBody map[string]interface{}, key string) *string {
	if val, ok := requestBody[key]; ok {
		if strVal, ok := val.(string); ok && strVal != "" {
			return &strVal
		}
	}
	return nil
}

func nullableString(str string) *string {
	if str == "" {
		return nil
	}
	return &str
}
