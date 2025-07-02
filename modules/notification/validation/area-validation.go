package validation

import (
	"byu-crm-service/helper"

	"github.com/go-playground/validator/v10"
)

// Validator instance
var validate = validator.New()

// CreateAreaRequest struct for creating a Area
type CreateAreaRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateAreaRequest struct {
	Name string `json:"name" validate:"required"`
}

// Mapping Validation Messages
var validationMessages = map[string]string{
	"name.required": "Nama Area harus diisi",
}

// ValidateCreate function to validate CreateAreaRequest
func ValidateCreate(req *CreateAreaRequest) map[string]string {
	return helper.ValidateStruct(validate, req, validationMessages)
}

// ValidateUpdate function to validate UpdateAreaRequest
func ValidateUpdate(req *UpdateAreaRequest) map[string]string {
	return helper.ValidateStruct(validate, req, validationMessages)
}
