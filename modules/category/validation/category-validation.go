package validation

import (
	"byu-crm-service/helper"

	"github.com/go-playground/validator/v10"
)

// Validator instance
var validate = validator.New()

// CreateCategoryRequest struct for creating a faculty
type CreateCategoryRequest struct {
	Name       string  `json:"name" validate:"required"`
	ModuleType *string `json:"module_type" validate:"required"`
}

type UpdateCategoryRequest struct {
	Name       string  `json:"name" validate:"required"`
	ModuleType *string `json:"module_type" validate:"required"`
}

// Mapping Validation Messages
var validationMessages = map[string]string{
	"name.required":        "Nama Tipe harus diisi",
	"module_type.required": "Tipe modul harus diisi",
}

// ValidateCreate function to validate CreateCategoryRequest
func ValidateCreate(req *CreateCategoryRequest) map[string]string {
	return helper.ValidateStruct(validate, req, validationMessages)
}

func ValidateUpdate(req *UpdateCategoryRequest) map[string]string {
	return helper.ValidateStruct(validate, req, validationMessages)
}
