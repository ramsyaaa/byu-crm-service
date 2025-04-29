package validation

import (
	"byu-crm-service/helper"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidateRequest struct {
	ContactName string  `json:"contact_name" validate:"required"`
	PhoneNumber string  `json:"phone_number" validate:"required"`
	Position    string  `json:"position" validate:"required"`
	Birthday    *string `json:"birthday"`

	AccountID []string `json:"account_id"`
	Category  []string `json:"category" validate:"validate_category_url"`
	Url       []string `json:"url"`
}

var validationMessages = map[string]string{
	"contact_name.required":          "Nama kontak wajib diisi",
	"phone_number.required":          "Nomor HP wajib diisi",
	"position.required":              "Posisi wajib diisi",
	"category.validate_category_url": "Kategori dan URL wajib diisi",
}

var validate = validator.New()

func init() {
	validate.RegisterValidation("validate_category_url", validateCategoryAndUrl)
}

func ValidateCreate(req *ValidateRequest) map[string]string {
	return helper.ValidateStruct(validate, req, validationMessages)
}

func validateCategoryAndUrl(fl validator.FieldLevel) bool {
	parent := fl.Parent()

	categoryField := parent.FieldByName("Category")
	urlField := parent.FieldByName("Url")

	if !categoryField.IsValid() || !urlField.IsValid() {
		return true
	}

	categories, ok1 := categoryField.Interface().([]string)
	urls, ok2 := urlField.Interface().([]string)

	if !ok1 || !ok2 {
		return false
	}

	if len(categories) != len(urls) {
		return false
	}

	for i := range categories {
		if strings.TrimSpace(categories[i]) == "" || strings.TrimSpace(urls[i]) == "" {
			return false
		}
	}

	return true
}
