package validation

import (
	"byu-crm-service/helper"

	"github.com/go-playground/validator/v10"
)

// Validator instance
var validate = validator.New()

// UpdateProfileRequest struct for creating a faculty
type UpdateProfileRequest struct {
	OldPassword     string `json:"old_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8"`
}

// Mapping Validation Messages
var validationMessages = map[string]string{
	"OldPassword.required":     "Password lama harus diisi",
	"NewPassword.required":     "Password baru harus diisi",
	"NewPassword.min":          "Password baru minimal 8 karakter",
	"ConfirmPassword.required": "Konfirmasi password harus diisi",
	"ConfirmPassword.min":      "Konfirmasi password minimal 8 karakter",
}

// ValidateUpdate function to validate UpdateProfileRequest
func ValidateUpdate(req *UpdateProfileRequest) map[string]string {
	err := validate.Struct(req)
	if err != nil {
		return helper.ErrorValidationFormat(err, validationMessages)
	}
	return nil
}
