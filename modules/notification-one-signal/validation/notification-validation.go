package validation

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type SubscribeNotificationRequest struct {
	Type           string `json:"type" validate:"required"`
	SubscriptionID string `json:"subscription_id" validate:"required"`
}

var validationMessages = map[string]string{
	"type.required":            "Tipe harus diisi",
	"subscription_id.required": "Subscription ID harus diisi",
}

func ValidateCreate(req *SubscribeNotificationRequest) map[string]string {
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
