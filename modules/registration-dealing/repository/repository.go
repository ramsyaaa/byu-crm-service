package repository

import (
	"byu-crm-service/modules/registration-dealing/response"
)

type RegistrationDealingRepository interface {
	GetAllRegistrationDealings(limit int, paginate bool, page int, filters map[string]string, accountID int, eventName string) ([]response.RegistrationDealingResponse, int64, error)
	GetAllRegistrationDealingGrouped(limit int, paginate bool, page int, filters map[string]string) ([]map[string]interface{}, int64, error)
	FindByRegistrationDealingID(id uint) (*response.RegistrationDealingResponse, error)
	CreateRegistrationDealing(requestBody map[string]string, userID *int) (*response.RegistrationDealingResponse, error)
	FindByPhoneNumber(phone_number string) (*response.RegistrationDealingResponse, error)
}
