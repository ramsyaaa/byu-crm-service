package validation

import (
	"byu-crm-service/helper"

	"github.com/go-playground/validator/v10"
)

// Validator instance
var validate = validator.New()

// CreateConstantDataRequest struct for creating a Area
type CreateConstantDataRequest struct {
	Value      string `json:"value" validate:"required"`
	Label      string `json:"label" validate:"required"`
	Type       string `json:"type" validate:"required"`
	OtherGroup string `json:"other_group"`
}

// Mapping Validation Messages
var validationMessages = map[string]string{
	"value.required": "Value harus diisi",
	"label.required": "Label harus diisi",
	"type.required":  "Type harus diisi",
}

// ValidateCreate function to validate CreateConstantDataRequest
func ValidateCreate(req *CreateConstantDataRequest) map[string]string {
	return helper.ValidateStruct(validate, req, validationMessages)
}
