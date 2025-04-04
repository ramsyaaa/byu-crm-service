package validation

import (
	"byu-crm-service/helper"

	"github.com/go-playground/validator/v10"
)

// Validator instance
var validate = validator.New()

// CreateAbsenceUserRequest struct for creating a absence
type CreateAbsenceUserRequest struct {
	Type        string `json:"type" validate:"required"`
	ActionType  string `json:"action_type"`
	Description string `json:"description"`
	Longitude   string `json:"longitude" validate:"required"`
	Latitude    string `json:"latitude" validate:"required"`
}

var validationMessages = map[string]string{
	"Type.required":       "Tipe harus diisi",
	"ActionType.required": "Tipe aksi harus diisi",
	"Longitude.required":  "Longitude harus diisi",
	"Latitude.required":   "Latitude harus diisi",
}

// ValidateCreate function to validate CreateAbsenceUserRequest
func ValidateCreate(req *CreateAbsenceUserRequest) map[string]string {
	err := validate.Struct(req)
	if err != nil {
		return helper.ErrorValidationFormat(err, validationMessages)
	}
	return nil
}
