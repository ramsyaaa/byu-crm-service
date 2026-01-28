package http

import (
	"byu-crm-service/helper"
	"byu-crm-service/modules/performance-digipos/service"
	"context"
	"encoding/csv"
	"fmt"
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
	}()
	response := helper.APIResponse("File has been successfully received and is being processed.", fiber.StatusOK, "success", nil)
	return c.Status(fiber.StatusOK).JSON(response)
}
