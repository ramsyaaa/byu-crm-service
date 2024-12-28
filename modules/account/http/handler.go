package http

import (
	"context"
	"encoding/csv"
	"os"
	"time"

	"byu-crm-service/modules/account/service"
	"byu-crm-service/modules/account/validation"

	"github.com/gofiber/fiber/v2"
)

type AccountHandler struct {
	service service.AccountService
}

func NewAccountHandler(service service.AccountService) *AccountHandler {
	return &AccountHandler{service: service}
}

func (h *AccountHandler) Import(c *fiber.Ctx) error {
	// Validate the uploaded file
	if err := validation.ValidateAccountRequest(c); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Save file temporarily
	file, _ := c.FormFile("file_csv")
	tempPath := "./temp/" + file.Filename

	// Ensure the temp directory exists
	if _, err := os.Stat("./temp/"); os.IsNotExist(err) {
		if err := os.Mkdir("./temp/", os.ModePerm); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create temp directory"})
		}
	}

	if err := c.SaveFile(file, tempPath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save file"})
	}
	defer os.Remove(tempPath)

	// Process file with timeout
	_, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	f, err := os.Open(tempPath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to open file"})
	}
	defer f.Close()

	reader := csv.NewReader(f)
	rows, _ := reader.ReadAll()

	for i, row := range rows {
		if i == 0 {
			continue // Skip header
		}
		if err := h.service.ProcessAccount(row); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	return c.JSON(fiber.Map{"message": "File processed successfully"})
}
