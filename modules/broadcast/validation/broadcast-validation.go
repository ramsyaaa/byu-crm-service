package validation

import (
	"byu-crm-service/helper"

	"github.com/go-playground/validator/v10"
)

type ValidateRequest struct {
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Type        string    `json:"type" validate:"required"`
	UserID      *[]string `json:"user_id"`
	RoleID      *[]string `json:"role_id"`
}

var validationMessages = map[string]string{
	"title.required":       "Judul wajib diisi",
	"description.required": "Deskripsi wajib diisi",
	"type.required":        "Tipe wajib diisi",
}

var validate = validator.New()

func ValidateCreate(req *ValidateRequest) map[string]string {
	return helper.ValidateStruct(validate, req, validationMessages)
}

func ValidateByUser(req *ValidateRequest) map[string]string {
	errors := make(map[string]string)
	if req.UserID == nil || len(*req.UserID) == 0 {
		errors["user_id"] = "User ID tidak boleh kosong."
	}

	if len(errors) == 0 {
		return nil
	}
	return errors
}

func ValidateByRole(req *ValidateRequest) map[string]string {
	errors := make(map[string]string)
	if req.RoleID == nil || len(*req.RoleID) == 0 {
		errors["role_id"] = "Role ID tidak boleh kosong."
	}

	if len(errors) == 0 {
		return nil
	}
	return errors
}
