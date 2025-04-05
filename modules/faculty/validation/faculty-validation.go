package validation

import (
	"byu-crm-service/helper"

	"github.com/go-playground/validator/v10"
)

// Validator instance
var validate = validator.New()

// CreateFacultyRequest struct for creating a faculty
type CreateFacultyRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateFacultyRequest struct {
	Name string `json:"name" validate:"required"`
}

// Mapping Validation Messages
var validationMessages = map[string]string{
	"Name.required": "Nama Fakultas harus diisi",
}

// ValidateCreate function to validate CreateFacultyRequest
func ValidateCreate(req *CreateFacultyRequest) map[string]string {
	err := validate.Struct(req)
	if err != nil {
		return helper.ErrorValidationFormat(err, validationMessages)
	}
	return nil
}

func ValidateUpdate(req *UpdateFacultyRequest) map[string]string {
	err := validate.Struct(req)
	if err != nil {
		return helper.ErrorValidationFormat(err, validationMessages)
	}
	return nil
}
