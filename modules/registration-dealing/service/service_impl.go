package service

import (
	"byu-crm-service/modules/registration-dealing/repository"
	"byu-crm-service/modules/registration-dealing/response"
	"fmt"
)

type registrationDealingService struct {
	repo repository.RegistrationDealingRepository
}

func NewRegistrationDealingService(repo repository.RegistrationDealingRepository) RegistrationDealingService {
	return &registrationDealingService{repo: repo}
}

func (s *registrationDealingService) GetAllRegistrationDealings(limit int, paginate bool, page int, filters map[string]string, accountID int, eventName string) ([]response.RegistrationDealingResponse, int64, error) {
	return s.repo.GetAllRegistrationDealings(limit, paginate, page, filters, accountID, eventName)
}

func (s *registrationDealingService) FindByRegistrationDealingID(id uint) (*response.RegistrationDealingResponse, error) {
	registrationDealing, err := s.repo.FindByRegistrationDealingID(id)
	if err != nil {
		return nil, err
	}

	return registrationDealing, nil
}

func (s *registrationDealingService) CreateRegistrationDealing(requestBody map[string]interface{}, userID int) (*response.RegistrationDealingResponse, error) {
	// Use getStringValue to safely handle nil values and type conversions
	registrationDealingData := map[string]string{
		"phone_number":          getStringValue(requestBody["phone_number"]),
		"account_id":            getStringValue(requestBody["account_id"]),
		"customer_name":         getStringValue(requestBody["customer_name"]),
		"event_name":            getStringValue(requestBody["event_name"]),
		"registration_evidence": getStringValue(requestBody["registration_evidence"]),
		"whatsapp_number":       getStringValue(requestBody["whatsapp_number"]),
		"class":                 getStringValue(requestBody["class"]),
		"email":                 getStringValue(requestBody["email"]),
		"user_id":               getStringValue(userID),
		"school_type":           getStringValue(requestBody["school_type"]),
	}

	registrationDealing, err := s.repo.CreateRegistrationDealing(registrationDealingData, userID)
	if err != nil {
		return nil, err
	}

	return registrationDealing, nil
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
