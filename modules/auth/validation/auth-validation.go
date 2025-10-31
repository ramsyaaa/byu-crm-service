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

type ImpersonateRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// Google OAuth Callback Request
type GoogleCallbackRequest struct {
	Code string `json:"code" validate:"required"`
}

// Mapping Pesan Validasi
var validationMessages = map[string]string{
	"Email.required":    "Email is required",
	"Email.email":       "Invalid email format",
	"Password.required": "Password is required",
	"Password.min":      "Password must be at least 6 characters long",
	"Code.required":     "Authorization code is required",
}

// Fungsi Validasi Request
func ValidateLogin(req *LoginRequest) map[string]string {
	err := validate.Struct(req)
	if err != nil {
		return helper.ErrorValidationFormat(err, validationMessages)
	}
	return nil
}

func ValidateImpersonate(req *ImpersonateRequest) map[string]string {
	err := validate.Struct(req)
	if err != nil {
		return helper.ErrorValidationFormat(err, validationMessages)
	}
	return nil
}

// Fungsi Validasi Google Callback Request
func ValidateGoogleCallback(req *GoogleCallbackRequest) map[string]string {
	err := validate.Struct(req)
	if err != nil {
		return helper.ErrorValidationFormat(err, validationMessages)
	}
	return nil
}
