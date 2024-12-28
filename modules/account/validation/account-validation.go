package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UploadRequest struct {
	FileCSV string `form:"file_csv" validate:"required,file_extension=csv"`
}

func ValidateAccountRequest(c *fiber.Ctx) error {
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
