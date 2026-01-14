package http

import (
	"bytes"
	"byu-crm-service/modules/performance-digipos/service"
	"byu-crm-service/modules/performance-digipos/validation"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type PerformanceDigiposHandler struct {
	service service.PerformanceDigiposService
}

func NewPerformanceDigiposHandler(service service.PerformanceDigiposService) *PerformanceDigiposHandler {
	return &PerformanceDigiposHandler{service: service}
}

func (h *PerformanceDigiposHandler) Import(c *fiber.Ctx) error {
	// Validate the uploaded file
	if err := validation.ValidatePerformanceDigiposRequest(c); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Save file temporarily
	file, err := c.FormFile("file_csv")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "File tidak ditemukan"})
	}

	filename := file.Filename
	ext := strings.ToLower(filepath.Ext(filename))

	// Validasi ekstensi
	allowed := map[string]bool{
		".csv": true,
	}

	if !allowed[ext] {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "Format file harus CSV atau Excel"})
	}

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
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID is required"})
	}

	// Menghitung jumlah total baris untuk estimasi durasi
	totalRows, err := countCSVRows(tempPath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to count CSV rows"})
	}

	// Asumsi setiap baris membutuhkan 0.5 detik untuk diproses
	processingTimePerRow := 0.5
	estimatedDuration := time.Duration(float64(totalRows) * processingTimePerRow * float64(time.Second))

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
			if err := h.service.ProcessPerformanceDigipos(row); err != nil {
				fmt.Println("Error processing row:", err)
				return
			}
		}

		// Send notification
		fmt.Println("Sending notification...")
		notificationURL := os.Getenv("NOTIFICATION_URL") + "/api/notification/create"
		payload := map[string]interface{}{
			"model":    "App\\Models\\Performance",
			"model_id": 0, // Replace with actual model ID if needed
			"user_id":  userID,
			"data": map[string]string{
				"title":        "Import Performance Digipos",
				"description":  "Import Performance Digipos",
				"callback_url": "/performances-digiposId",
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

	return c.JSON(fiber.Map{
		"message":           "File upload successful, processing in background",
		"estimated_seconds": estimatedDuration.Seconds(),
	})
}

// countCSVRows menghitung jumlah total baris dalam file CSV
func countCSVRows(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	totalRows := 0

	// Baca setiap baris dan hitung jumlah totalnya
	for {
		_, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return 0, err
		}
		totalRows++
	}

	return totalRows, nil
}
