package validation

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type CreateAbsenceUserRequest struct {
	Type        string `json:"type" validate:"required"`
	ActionType  string `json:"action_type"`
	Description string `json:"description"`
	Longitude   string `json:"longitude" validate:"required"`
	Latitude    string `json:"latitude" validate:"required"`
}

var validationMessages = map[string]string{
	"type.required":        "Tipe harus diisi",
	"action_type.required": "Tipe aksi harus diisi",
	"longitude.required":   "Longitude harus diisi",
	"latitude.required":    "Latitude harus diisi",
}

func ValidateCreate(req *CreateAbsenceUserRequest) map[string]string {
	err := validate.Struct(req)
	if err == nil {
		return nil
	}

	errors := make(map[string]string)
	validationErrors := err.(validator.ValidationErrors)
	ref := reflect.TypeOf(*req)

	for _, e := range validationErrors {
		// Ambil nama json tag dari field
		field, _ := ref.FieldByName(e.StructField())
		jsonTag := field.Tag.Get("json")
		jsonKey := strings.Split(jsonTag, ",")[0]

		// Ambil pesan error dari map validationMessages
		key := jsonKey + "." + e.Tag()
		msg, found := validationMessages[key]
		if !found {
			msg = e.Error() // fallback kalau gak ada
		}
		errors[jsonKey] = msg // Gunakan json key sebagai key utama
	}

	return errors
}
