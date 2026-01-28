package http

import (
	"byu-crm-service/helper"
	"byu-crm-service/modules/performance-indiana/service"
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

type PerformanceIndianaHandler struct {
	service service.PerformanceIndianaService
}

func NewPerformanceIndianaHandler(service service.PerformanceIndianaService) *PerformanceIndianaHandler {
	return &PerformanceIndianaHandler{service: service}
}

func (h *PerformanceIndianaHandler) Import(c *fiber.Ctx) error {
	fileBase64 := c.FormValue("file_csv")
	monthYear := c.FormValue("month") // format: YYYY-MM

	if fileBase64 == "" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "File wajib diisi"})
	}

	if monthYear == "" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "Bulan wajib diisi"})
	}

	parsedTime, err := time.Parse("2006-01", monthYear)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "Format bulan tidak valid (YYYY-MM)"})
	}

	month := uint(parsedTime.Month())
	year := uint(parsedTime.Year())

	result, err := helper.DecodeBase64Excel(fileBase64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	tempPath := result.Path

	// 🚀 Async process
	go func(importMonth uint, importYear uint) {
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
				if err := h.service.ProcessPerformanceIndiana(
					row,
					importMonth,
					importYear,
				); err != nil {
					fmt.Println("Error processing row:", err)
				}
			}
		}
	}(month, year)

	response := helper.APIResponse(
		"File has been successfully received and is being processed.",
		fiber.StatusOK,
		"success",
		nil,
	)

	return c.Status(fiber.StatusOK).JSON(response)
}
