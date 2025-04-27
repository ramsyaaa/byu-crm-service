package validation

import (
	"byu-crm-service/helper"

	"github.com/go-playground/validator/v10"
)

// Validator instance
var validate = validator.New()

// CreateTypeRequest struct for creating a faculty
type CreateTypeRequest struct {
	Name        string  `json:"name" validate:"required"`
	CategoryID  string  `json:"category_id"`
	ModuleType  *string `json:"module_type" validate:"required"`
	Description string  `json:"description"`
}

type UpdateTypeRequest struct {
	Name        string  `json:"name" validate:"required"`
	CategoryID  string  `json:"category_id"`
	ModuleType  *string `json:"module_type" validate:"required"`
	Description string  `json:"description"`
}

// Mapping Validation Messages
var validationMessages = map[string]string{
	"Name.required":       "Nama Tipe harus diisi",
	"ModuleType.required": "Tipe modul harus diisi",
}

// ValidateCreate function to validate CreateTypeRequest
func ValidateCreate(req *CreateTypeRequest) map[string]string {
	return helper.ValidateStruct(validate, req, validationMessages)
}

func ValidateUpdate(req *UpdateTypeRequest) map[string]string {
	return helper.ValidateStruct(validate, req, validationMessages)
}
