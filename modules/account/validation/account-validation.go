package validation

import (
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UploadRequest struct {
	UserID string `form:"user_id" validate:"required"`
}

func ValidateAccountRequest(c *fiber.Ctx) error {
	// Check if file exists in the request
	file, err := c.FormFile("file_csv")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File is required",
		})
	}

	// Validate file extension
	if !validateFileExtension(file, "csv") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Only CSV files are allowed",
		})
	}

	// Validate user_id
	var request UploadRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Next()
}

// Helper function to validate file extension
func validateFileExtension(file *multipart.FileHeader, allowedExt string) bool {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	return ext == "."+allowedExt
}
