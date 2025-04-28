package validation

import (
	"byu-crm-service/helper"

	"github.com/go-playground/validator/v10"
)

type ValidateRequest struct {
	ContactName string  `json:"contact_name" validate:"required"`
	PhoneNumber string  `json:"phone_number" validate:"required"`
	Position    string  `json:"position" validate:"required"`
	Birthday    *string `json:"birthday"`

	AccountID []string `json:"account_id"`
}

var validationMessages = map[string]string{
	"contact_name.required": "Nama kontak wajib diisi",
	"phone_number.required": "Nomor HP wajib diisi",
	"position.required":     "Posisi wajib diisi",
}

var validate = validator.New()

func ValidateCreate(req *ValidateRequest) map[string]string {
	return helper.ValidateStruct(validate, req, validationMessages)
}
