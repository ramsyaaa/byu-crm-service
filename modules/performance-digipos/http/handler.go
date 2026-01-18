package http

import (
	"byu-crm-service/helper"
	"byu-crm-service/modules/performance-digipos/service"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
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
	fileBase64 := c.FormValue("file_csv")

	if fileBase64 == "" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "File wajib diisi"})
	}

	result, err := helper.DecodeBase64Excel(fileBase64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	tempPath := result.Path

	totalRows, err := countCSVRows(tempPath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": "Gagal membaca file CSV"})
	}

	processingTimePerRow := 0.5
	estimatedDuration := time.Duration(
		float64(totalRows) * processingTimePerRow * float64(time.Second),
	)

	go func() {
		defer os.Remove(tempPath)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		file, err := os.Open(tempPath)
		if err != nil {
			fmt.Println("Gagal membuka file:", err)
			return
		}
		defer file.Close()

		reader := csv.NewReader(file)
		rows, err := reader.ReadAll()
		if err != nil {
			fmt.Println("Gagal membaca CSV:", err)
			return
		}

		for i, row := range rows {
			if i == 0 {
				continue // skip header
			}

			select {
			case <-ctx.Done():
				fmt.Println("Proses import timeout")
				return
			default:
				if err := h.service.ProcessPerformanceDigipos(row); err != nil {
					fmt.Println("Error processing row:", err)
				}
			}
		}

		// =====================
		// KIRIM NOTIFIKASI
		// =====================
		// notificationURL := os.Getenv("NOTIFICATION_URL") + "/api/notification/create"

		// payload := map[string]interface{}{
		// 	"model":    "App\\Models\\Performance",
		// 	"model_id": 0,
		// 	"user_id":  userID,
		// 	"data": map[string]string{
		// 		"title":        "Import Performance Digipos",
		// 		"description":  "Import Performance Digipos berhasil",
		// 		"callback_url": "/performances-digipos",
		// 	},
		// }

		// payloadBytes, _ := json.Marshal(payload)

		// resp, err := http.Post(
		// 	notificationURL,
		// 	"application/json",
		// 	bytes.NewReader(payloadBytes),
		// )
		// if err != nil {
		// 	fmt.Println("Gagal kirim notifikasi:", err)
		// 	return
		// }
		// defer resp.Body.Close()

		fmt.Println("Import Performance Digipos selesai")

	}()

	// =====================
	// RESPONSE USER
	// =====================
	return c.JSON(fiber.Map{
		"message":           "File berhasil diterima dan sedang diproses",
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
