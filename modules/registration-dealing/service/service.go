package service

import (
	"byu-crm-service/modules/registration-dealing/response"
)

type RegistrationDealingService interface {
	GetAllRegistrationDealings(limit int, paginate bool, page int, filters map[string]string, accountID int, event_name string) ([]response.RegistrationDealingResponse, int64, error)
	FindByRegistrationDealingID(id uint) (*response.RegistrationDealingResponse, error)
	CreateRegistrationDealing(requestBody map[string]interface{}, userID int) (*response.RegistrationDealingResponse, error)
}
