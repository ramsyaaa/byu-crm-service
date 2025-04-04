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
	Description string `json:"description" validate:"required"`
	Longitude   string `json:"longitude" validate:"required"`
	Latitude    string `json:"latitude" validate:"required"`
}

var validationMessages = map[string]string{
	"Type.required":        "Tipe harus diisi",
	"Description.required": "Deskripsi harus diisi",
	"Longitude.required":   "Longitude harus diisi",
	"Latitude.required":    "Latitude harus diisi",
}

// ValidateCreate function to validate CreateAbsenceUserRequest
func ValidateCreate(req *CreateAbsenceUserRequest) map[string]string {
	err := validate.Struct(req)
	if err != nil {
		return helper.ErrorValidationFormat(err, validationMessages)
	}
	return nil
}
