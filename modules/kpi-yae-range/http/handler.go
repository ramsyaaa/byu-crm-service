package http

import (
	"byu-crm-service/helper"
	"byu-crm-service/modules/kpi-yae-range/service"
	"byu-crm-service/modules/kpi-yae-range/validation"
	kpiYaeService "byu-crm-service/modules/kpi-yae/service"
	visitHistoryService "byu-crm-service/modules/visit-history/service"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type KpiYaeRangeHandler struct {
	service             service.KpiYaeRangeService
	kpiYaeService       kpiYaeService.KpiYaeService
	visitHistoryService visitHistoryService.VisitHistoryService
}

func NewKpiYaeRangeHandler(service service.KpiYaeRangeService, kpiYaeService kpiYaeService.KpiYaeService, visitHistoryService visitHistoryService.VisitHistoryService) *KpiYaeRangeHandler {
	return &KpiYaeRangeHandler{service: service, kpiYaeService: kpiYaeService, visitHistoryService: visitHistoryService}
}

func (h *KpiYaeRangeHandler) GetCurrentKpiYaeRanges(c *fiber.Ctx) error {
	// Get the current KpiYaeRanges from the service
	now := time.Now()
	month := uint(now.Month())
	year := uint(now.Year())

	kpiYaeRanges, err := h.service.GetKpiYaeRangeByDate(month, year)
	if err != nil {
		response := helper.APIResponse("Failed to fetch KPI", fiber.StatusInternalServerError, "error", nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Return response
	responseData := map[string]interface{}{
		"kpi_yae_ranges": kpiYaeRanges,
	}

	response := helper.APIResponse("Get Current KpiYaeRanges Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *KpiYaeRangeHandler) CreateKpiYaeRange(c *fiber.Ctx) error {
	req := new(validation.CreateKpiYaeRangeRequest)
	if err := c.BodyParser(req); err != nil {
		response := helper.APIResponse("Invalid request", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Request Validation
	errors := validation.ValidateCreate(req)
	if errors != nil {
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if len(req.Name) != len(req.Target) {
		errors := map[string]string{
			"Name":   "Jumlah nama KPI tidak sama dengan jumlah target KPI",
			"Target": "Jumlah target KPI tidak sama dengan jumlah nama KPI",
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	var invalidNames []string

	for _, name := range req.Name {
		kpi, err := h.kpiYaeService.GetKpiYaeByName(name)
		if err != nil || kpi == nil {
			invalidNames = append(invalidNames, name)
		}
	}

	if len(invalidNames) > 0 {
		errors := map[string]string{
			"Name": "KPI tidak ditemukan",
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// // Create a new KpiYaeRange using the service
	kpiYaeRange, err := h.service.CreateKpiYaeRange(req)
	if err != nil {
		response := helper.APIResponse("Failed to create KPI", fiber.StatusInternalServerError, "error", nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// // Return response
	responseData := map[string]interface{}{
		"kpi_yae_range": kpiYaeRange,
	}

	response := helper.APIResponse("Create KPI YAE Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

type Item struct {
	Name   string `json:"name"`
	Target string `json:"target"`
}

type UserPerformance struct {
	Name       string `json:"name"`
	Target     string `json:"target"`
	Actual     string `json:"actual"`
	Percentage string `json:"percentage"`
}

func (h *KpiYaeRangeHandler) GetPerformanceUser(c *fiber.Ctx) error {
	now := time.Now()
	month := uint(now.Month())
	year := uint(now.Year())
	userID := c.Locals("user_id").(int)

	kpi_lists, err := h.service.GetKpiYaeRangeByDate(month, year)
	if err != nil {
		response := helper.APIResponse("Failed to fetch KPI", fiber.StatusInternalServerError, "error", nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	var items []Item
	var performances []UserPerformance

	err = json.Unmarshal([]byte(kpi_lists.Target), &items)
	if err != nil {
		response := helper.APIResponse("Error data KPI", fiber.StatusInternalServerError, "error", nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Looping KPI
	for _, item := range items {
		if item.Name == "Visit" {
			visitActual, err := h.visitHistoryService.CountVisitHistory(userID, month, year, "")
			if err != nil {
				response := helper.APIResponse("Counting Visit Error", fiber.StatusInternalServerError, "error", nil)
				return c.Status(fiber.StatusInternalServerError).JSON(response)
			}

			target, _ := strconv.Atoi(item.Target) // konversi target ke int
			percentage := 0
			if target > 0 {
				percentage = (visitActual * 100) / target
			}

			performances = append(performances, UserPerformance{
				Name:       item.Name,
				Target:     item.Target,
				Actual:     strconv.Itoa(visitActual),
				Percentage: fmt.Sprintf("%d%%", percentage),
			})
		} else if item.Name == "Presentasi Demo" {
			PresentationActual, err := h.visitHistoryService.CountVisitHistory(userID, month, year, item.Name)
			if err != nil {
				response := helper.APIResponse("Counting Visit Error", fiber.StatusInternalServerError, "error", nil)
				return c.Status(fiber.StatusInternalServerError).JSON(response)
			}

			target, _ := strconv.Atoi(item.Target) // convertion to int
			percentage := 0
			if target > 0 {
				percentage = (PresentationActual * 100) / target
			}

			performances = append(performances, UserPerformance{
				Name:       item.Name,
				Target:     item.Target,
				Actual:     strconv.Itoa(PresentationActual),
				Percentage: fmt.Sprintf("%d%%", percentage),
			})
		}
	}
	// Return response
	responseData := map[string]interface{}{
		"performances": performances,
	}

	response := helper.APIResponse("Get Performance Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}
