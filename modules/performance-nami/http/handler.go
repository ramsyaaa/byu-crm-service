package http

import (
	"byu-crm-service/modules/performance-nami/service"
	"byu-crm-service/modules/performance-nami/validation"
	"context"
	"encoding/csv"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Struct untuk performance
type PerformanceNami struct {
	Periode            string
	PeriodeDate        *time.Time
	EventID            string
	PoiID              string
	PoiName            string
	PoiType            string
	EventName          string
	EventType          string
	EventLocationType  string
	SalesType          string
	SalesType2         string
	CityID             *uint
	SerialNumberMSISDN string
	ScanType           string
	ActiveMSISDN       string
	ActiveDate         *time.Time
	ActiveCity         string
	Validation         string
	ValidKPI           bool
	Rev                string
	SaDate             *time.Time
	SoDate             *time.Time
	NewImei            string
	SkulIDDate         *time.Time
	AgentID            string
	UserID             string
	UserName           string
	UserType           string
	UserSubType        string
	ScanDate           *time.Time
	Plan               string
	TopStatus          bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

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
		if err := h.service.ProcessPerformanceNami(row); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	return c.JSON(fiber.Map{"message": "File processed successfully"})
}
