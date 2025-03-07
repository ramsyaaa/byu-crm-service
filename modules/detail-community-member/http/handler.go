package http

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"byu-crm-service/modules/detail-community-member/service"
	"byu-crm-service/modules/detail-community-member/validation"

	"github.com/gofiber/fiber/v2"
)

type DetailCommunityMemberHandler struct {
	service service.DetailCommunityMemberService
}

func NewDetailCommunityMemberHandler(service service.DetailCommunityMemberService) *DetailCommunityMemberHandler {
	return &DetailCommunityMemberHandler{service: service}
}

func (h *DetailCommunityMemberHandler) Import(c *fiber.Ctx) error {
	// Validate the uploaded file
	if err := validation.ValidateDetailCommunityMemberRequest(c); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Save file temporarily
	file, _ := c.FormFile("file_csv_community")
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
	accountIDStr := c.FormValue("account_id")
	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid account_id"})
	}

	go func() {
		defer os.Remove(tempPath)

		// Process file with timeout
		_, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		f, err := os.Open(tempPath)
		if err != nil {
			fmt.Println("Failed to open file:", err)
			return
		}
		defer f.Close()

		reader := csv.NewReader(f)
		rows, _ := reader.ReadAll()
		uploadedDate := time.Now()

		for i, row := range rows {
			if i == 0 {
				continue // Skip header
			}
			if err := h.service.ProcessData(row, uint(accountID), uploadedDate); err != nil {
				fmt.Println("Error processing row:", err)
				return
			}
		}
	}()

	return c.JSON(fiber.Map{"message": "File processed successfully"})
}
