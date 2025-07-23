package validation

import (
	"byu-crm-service/helper"

	"github.com/go-playground/validator/v10"
)

type ValidateRequest struct {
	Title         string    `json:"title" validate:"required"`
	Description   string    `json:"description" validate:"required"`
	Type          string    `json:"type" validate:"required"`
	UserID        *[]string `json:"user_id"`
	RoleID        *[]string `json:"role_id"`
	TerritoryType *string   `json:"territory_type"`
	TerritoryID   *[]string `json:"territory_id"`
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
		errors["user_id"] = "User tidak boleh kosong."
	}

	if len(errors) == 0 {
		return nil
	}
	return errors
}

func ValidateByRole(req *ValidateRequest) map[string]string {
	errors := make(map[string]string)
	if req.RoleID == nil || len(*req.RoleID) == 0 {
		errors["role_id"] = "Role tidak boleh kosong."
	}

	if len(errors) == 0 {
		return nil
	}
	return errors
}

func ValidateByTerritory(req *ValidateRequest) map[string]string {
	errors := make(map[string]string)
	if req.TerritoryType == nil || *req.TerritoryType == "" {
		errors["territory_type"] = "Type territory tidak boleh kosong."
	}

	if *req.TerritoryType != "all" && *req.TerritoryType != "" {
		if req.TerritoryID == nil || len(*req.TerritoryID) == 0 {
			errors["territory_id"] = "Territory tidak boleh kosong."
		}
	}

	if *req.TerritoryType != "all" && *req.TerritoryType != "areas" && *req.TerritoryType != "regions" && *req.TerritoryType != "branches" && *req.TerritoryType != "clusters" {
		errors["territory_type"] = "Type territory tidak valid."
	}

	if len(errors) == 0 {
		return nil
	}
	return errors
}
