package service

import (
	"byu-crm-service/models"
	"byu-crm-service/modules/communication/repository"
	"encoding/json"
	"fmt"
)

type communicationService struct {
	repo repository.CommunicationRepository
}

func NewCommunicationService(repo repository.CommunicationRepository) CommunicationService {
	return &communicationService{repo: repo}
}

func (s *communicationService) GetAllCommunications(limit int, paginate bool, page int, filters map[string]string, accountID int) ([]models.Communication, int64, error) {
	return s.repo.GetAllCommunications(limit, paginate, page, filters, accountID)
}

func (s *communicationService) FindByCommunicationID(id uint) (*models.Communication, error) {
	communication, err := s.repo.FindByCommunicationID(id)
	if err != nil {
		return nil, err
	}

	return communication, nil
}

func (s *communicationService) CreateCommunication(requestBody map[string]interface{}, userID int) (*models.Communication, error) {
	// Use getStringValue to safely handle nil values and type conversions
	communicationData := map[string]string{
		"communication_type":    getStringValue(requestBody["communication_type"]),
		"note":                  getStringValue(requestBody["note"]),
		"account_id":            getStringValue(requestBody["account_id"]),
		"contact_id":            getStringValue(requestBody["contact_id"]),
		"created_by":            getStringValue(userID),
		"opportunity_id":        getStringValue(requestBody["opportunity_id"]),
		"status_communication":  getStringValue(requestBody["status_communication"]),
		"main_communication_id": getStringValue(""),
	}

	communication, err := s.repo.CreateCommunication(communicationData)
	if err != nil {
		return nil, err
	}

	return communication, nil
}

func (s *communicationService) UpdateCommunication(requestBody map[string]interface{}, userID int, communicationID int) (*models.Communication, error) {
	// Use getStringValue to safely handle nil values and type conversions

	previousCommunication, err := s.repo.FindByCommunicationID(uint(communicationID))
	if err != nil {
		return nil, err
	}
	if previousCommunication == nil {
		return nil, fmt.Errorf("communication with ID %d not found", communicationID)
	}

	var main_communication_id string
	if previousCommunication.MainCommunicationID != nil {
		main_communication_id = fmt.Sprintf("%d", *previousCommunication.MainCommunicationID)
	} else {
		main_communication_id = fmt.Sprintf("%d", previousCommunication.ID)
	}

	communicationData := map[string]string{
		"communication_type":    getStringValue(requestBody["communication_type"]),
		"note":                  getStringValue(requestBody["note"]),
		"account_id":            getStringValue(requestBody["account_id"]),
		"contact_id":            getStringValue(requestBody["contact_id"]),
		"created_by":            getStringValue(userID),
		"opportunity_id":        getStringValue(requestBody["opportunity_id"]),
		"status_communication":  getStringValue(requestBody["status_communication"]),
		"main_communication_id": main_communication_id,
	}

	communication, err := s.repo.CreateCommunication(communicationData)
	if err != nil {
		return nil, err
	}

	// Update the next communication ID in the previous communication record
	err = s.repo.UpdateFields(previousCommunication.ID, map[string]interface{}{"next_communication_id": communication.ID})

	if err != nil {
		return nil, err
	}

	return communication, nil
}

func SplitFields(jsonString string, keys []string) (map[string][]string, error) {
	// Cek jika JSON string kosong atau null
	if jsonString == "" || jsonString == "null" {
		// Jika kosong atau null, kembalikan array kosong untuk setiap key
		result := make(map[string][]string)
		for _, key := range keys {
			result[key] = []string{}
		}
		return result, nil
	}

	// Slice untuk menampung hasil decoding JSON
	var rawData []map[string]string

	// Decode JSON menjadi slice of map[string]string
	err := json.Unmarshal([]byte(jsonString), &rawData)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	// Map untuk menyimpan hasil pemisahan berdasarkan key
	result := make(map[string][]string)

	// Iterasi untuk setiap record dalam rawData
	for _, item := range rawData {
		// Iterasi setiap key yang ingin dipisahkan
		for _, key := range keys {
			// Cek apakah key ada dalam item, jika ada, tambahkan ke hasil
			if value, exists := item[key]; exists {
				result[key] = append(result[key], value)
			} else {
				// Jika key tidak ada, tambahkan array kosong
				result[key] = append(result[key], "")
			}
		}
	}

	return result, nil
}

func parseJSONStringToArray(jsonString string) ([]string, error) {
	// Jika jsonString kosong atau null, kembalikan array kosong
	if jsonString == "" || jsonString == "null" {
		return []string{}, nil
	}

	// Slice untuk menyimpan hasil decoding JSON
	var result []string

	// Decode JSON menjadi slice of string
	err := json.Unmarshal([]byte(jsonString), &result)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %v", err)
	}

	return result, nil
}

func getStringValue(val interface{}) string {
	if val == nil {
		return ""
	}
	if str, ok := val.(string); ok {
		return str
	}
	return fmt.Sprintf("%v", val)
}
