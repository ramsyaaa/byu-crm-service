package validation

import (
	"byu-crm-service/helper"

	"github.com/go-playground/validator/v10"
)

// Validator instance
var validate = validator.New()

// CreatePermissionRequest struct for creating a Area
type CreatePermissionRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdatePermissionRequest struct {
	Name string `json:"name" validate:"required"`
}

// Mapping Validation Messages
var validationMessages = map[string]string{
	"name.required": "Nama permission harus diisi",
}

// ValidateCreate function to validate CreatePermissionRequest
func ValidateCreate(req *CreatePermissionRequest) map[string]string {
	return helper.ValidateStruct(validate, req, validationMessages)
}

// ValidateUpdate function to validate UpdatePermissionRequest
func ValidateUpdate(req *UpdatePermissionRequest) map[string]string {
	return helper.ValidateStruct(validate, req, validationMessages)
}
