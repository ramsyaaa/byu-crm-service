package http

import (
	"bytes"
	"byu-crm-service/modules/performance-nami/service"
	"byu-crm-service/modules/performance-nami/validation"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

type PerformanceNamiHandler struct {
	service service.PerformanceNamiService
}

func NewPerformanceNamiHandler(service service.PerformanceNamiService) *PerformanceNamiHandler {
	return &PerformanceNamiHandler{service: service}
}

func (h *PerformanceNamiHandler) Import(c *fiber.Ctx) error {
	// Validate the uploaded file
	if err := validation.ValidatePerformanceNamiRequest(c); err != nil {
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

	// Retrieve user_id from the form
	userID := c.FormValue("user_id")

	// Respond immediately to the user
	go func() {
		defer os.Remove(tempPath) // Clean up the temporary file

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

		for i, row := range rows {
			if i == 0 {
				continue // Skip header
			}
			if err := h.service.ProcessPerformanceNami(row); err != nil {
				fmt.Println("Error processing row:", err)
				return
			}
		}

		// Send notification
		notificationURL := os.Getenv("NOTIFICATION_URL") + "/api/notification/create"
		payload := map[string]interface{}{
			"model":    "App\\Models\\PerformanceNami",
			"model_id": 0, // Replace with actual model ID if needed
			"user_id":  userID,
			"data": map[string]string{
				"title":        "Import Performance Nami",
				"description":  "Import Performance Nami",
				"callback_url": "/performances-nami",
			},
		}

		payloadBytes, _ := json.Marshal(payload)
		resp, err := http.Post(notificationURL, "application/json", bytes.NewReader(payloadBytes))
		if err != nil {
			fmt.Println("Failed to send notification:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println("Notification API responded with status:", resp.StatusCode)
		} else {
			var responseMap map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
				fmt.Println("Failed to decode response:", err)
				return
			}
			fmt.Println("Notification sent successfully:", responseMap["message"])
		}
	}()

	return c.JSON(fiber.Map{"message": "File upload successful, processing in background"})
}
