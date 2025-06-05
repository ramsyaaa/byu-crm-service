package validation

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

type UploadRequest struct {
	FileCSV string `form:"file_csv" validate:"required,file_extension=csv"`
	UserID  string `form:"user_id" validate:"required"`
}

type CreatePerformanceSkulIdRequest struct {
	UserName       string `json:"user_name" validate:"required"`
	IDSkulId       string `json:"id_skulid" validate:"required"`
	MSISDN         string `json:"msisdn" validate:"required"`
	RegisteredDate string `json:"registered_date" validate:"required"`
	Provider       string `json:"provider" validate:"required"`
	Batch          string `json:"batch" validate:"required"`
	UserType       string `json:"user_type" validate:"required,omitempty,oneof=Siswa Sekolah"`
}

var validationMessages = map[string]string{
	"user_name.required":       "Nama pengguna harus diisi",
	"id_skulid.required":       "ID Skul ID harus diisi",
	"msisdn.required":          "MSISDN harus diisi",
	"registered_date.required": "Tanggal pendaftaran harus diisi",
	"provider.required":        "Provider harus diisi",
	"batch.required":           "Angkatan harus diisi",
	"user_type.required":       "Tipe pengguna harus diisi",
	"user_type.oneof":          "Tipe pengguna harus salah satu dari: Siswa, Sekolah",
}

func ValidatePerformanceSkulIdRequest(c *fiber.Ctx) error {
	var request UploadRequest
	if err := c.BodyParser(&request); err != nil {
		return err
	}

	validate := validator.New()
	validate.RegisterValidation("file_extension", func(fl validator.FieldLevel) bool {
		return fl.Field().String() == "csv"
	})

	if err := validate.Struct(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Next()
}

func ValidateCreate(req *CreatePerformanceSkulIdRequest) map[string]string {
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
