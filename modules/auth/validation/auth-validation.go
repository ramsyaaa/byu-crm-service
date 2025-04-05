package validation

import (
	"byu-crm-service/helper"

	"github.com/go-playground/validator/v10"
)

// Validator instance
var validate = validator.New()

// Struktur Request dengan Tag Validasi
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// Mapping Pesan Validasi
var validationMessages = map[string]string{
	"Email.required":    "Email is required",
	"Email.email":       "Invalid email format",
	"Password.required": "Password is required",
	"Password.min":      "Password must be at least 6 characters long",
}

// Fungsi Validasi Request
func ValidateLogin(req *LoginRequest) map[string]string {
	err := validate.Struct(req)
	if err != nil {
		return helper.ErrorValidationFormat(err, validationMessages)
	}
	return nil
}
