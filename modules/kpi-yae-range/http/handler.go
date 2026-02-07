package http

import (
	"byu-crm-service/helper"
	"byu-crm-service/models"
	"byu-crm-service/modules/kpi-yae-range/service"
	"byu-crm-service/modules/kpi-yae-range/validation"
	kpiYaeService "byu-crm-service/modules/kpi-yae/service"
	performanceDigiposService "byu-crm-service/modules/performance-digipos/service"
	performanceIndianaService "byu-crm-service/modules/performance-indiana/service"
	userService "byu-crm-service/modules/user/service"
	visitHistoryService "byu-crm-service/modules/visit-history/service"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type KpiYaeRangeHandler struct {
	service                   service.KpiYaeRangeService
	kpiYaeService             kpiYaeService.KpiYaeService
	visitHistoryService       visitHistoryService.VisitHistoryService
	performanceDigiposService performanceDigiposService.PerformanceDigiposService
	userService               userService.UserService
	performanceIndianaService performanceIndianaService.PerformanceIndianaService
}

func NewKpiYaeRangeHandler(service service.KpiYaeRangeService, kpiYaeService kpiYaeService.KpiYaeService, visitHistoryService visitHistoryService.VisitHistoryService, performanceDigiposService performanceDigiposService.PerformanceDigiposService, userService userService.UserService, performanceIndianaService performanceIndianaService.PerformanceIndianaService) *KpiYaeRangeHandler {
	return &KpiYaeRangeHandler{service: service, kpiYaeService: kpiYaeService, visitHistoryService: visitHistoryService, performanceDigiposService: performanceDigiposService, userService: userService, performanceIndianaService: performanceIndianaService}
}

func (h *KpiYaeRangeHandler) GetCurrentKpiYaeRanges(c *fiber.Ctx) error {
	// Ambil query param dari URL jika ada
	monthParam := c.Query("month")
	yearParam := c.Query("year")

	// Default: bulan & tahun saat ini
	now := time.Now()
	month := uint(now.Month())
	year := uint(now.Year())

	// Jika parameter ada dan valid, override nilai default
	if monthParam != "" {
		if m, err := strconv.Atoi(monthParam); err == nil && m >= 1 && m <= 12 {
			month = uint(m)
		}
	}

	if yearParam != "" {
		if y, err := strconv.Atoi(yearParam); err == nil && y > 2000 {
			year = uint(y)
		}
	}

	// Ambil data dari service
	kpiYaeRanges, err := h.service.GetKpiYaeRangeByDate(month, year)
	if err != nil {
		if err.Error() == "record not found" {
			response := helper.APIResponse("KPI not found", fiber.StatusNotFound, "error", nil)
			return c.Status(fiber.StatusNotFound).JSON(response)
		}
		response := helper.APIResponse("Failed to fetch KPI", fiber.StatusInternalServerError, "error", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Return response
	responseData := map[string]interface{}{
		"kpi_yae_ranges": kpiYaeRanges,
	}
	response := helper.APIResponse("Get Current Kpi YAE Successfully", fiber.StatusOK, "success", responseData)
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

func (h *KpiYaeRangeHandler) GetPerformanceUser(c *fiber.Ctx) error {
	now := time.Now()
	month := uint(now.Month())
	year := uint(now.Year())
	userID := c.Locals("user_id").(int)

	if m := c.Query("month"); m != "" {
		if parsedMonth, err := strconv.Atoi(m); err == nil && parsedMonth >= 1 && parsedMonth <= 12 {
			month = uint(parsedMonth)
		}
	}

	// Cek jika parameter year dikirim
	if y := c.Query("year"); y != "" {
		if parsedYear, err := strconv.Atoi(y); err == nil && parsedYear > 0 {
			year = uint(parsedYear)
		}
	}

	paramUserID := c.Query("user_id")
	if paramUserID != "" {
		if parsedID, err := strconv.Atoi(paramUserID); err == nil {
			userID = parsedID
		}
	}

	if c.Query("all_user") == "1" {
		userID = 0
	}

	kpi_lists, err := h.service.GetKpiYaeRangeByDate(month, year)
	if err != nil {
		response := helper.APIResponse("Failed to fetch KPI", fiber.StatusInternalServerError, "error", nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	var items []Item
	var performances []models.UserPerformance

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

			performances = append(performances, models.UserPerformance{
				Name:       item.Name,
				Target:     item.Target,
				Actual:     strconv.Itoa(visitActual),
				Percentage: fmt.Sprintf("%d%%", percentage),
			})
		} else if item.Name == "Presentasi Demo" {
			PresentationActual, err := h.visitHistoryService.CountVisitHistory(userID, month, year, "presentasi_demo")
			if err != nil {
				response := helper.APIResponse("Counting Visit Error", fiber.StatusInternalServerError, "error", nil)
				return c.Status(fiber.StatusInternalServerError).JSON(response)
			}

			target, _ := strconv.Atoi(item.Target) // convertion to int
			percentage := 0
			if target > 0 {
				percentage = (PresentationActual * 100) / target
			}

			performances = append(performances, models.UserPerformance{
				Name:       item.Name,
				Target:     item.Target,
				Actual:     strconv.Itoa(PresentationActual),
				Percentage: fmt.Sprintf("%d%%", percentage),
			})
		} else if item.Name == "Digipos" {
			PresentationActual, err := h.performanceDigiposService.CountPerformanceByUserYaeCode(userID, month, year)
			if err != nil {
				fmt.Println(err.Error())
				response := helper.APIResponse("Counting Digipos Error", fiber.StatusInternalServerError, "error", nil)
				return c.Status(fiber.StatusInternalServerError).JSON(response)
			}

			target, _ := strconv.Atoi(item.Target) // convertion to int
			percentage := 0
			if target > 0 {
				percentage = (PresentationActual * 100) / target
			}

			performances = append(performances, models.UserPerformance{
				Name:       item.Name,
				Target:     item.Target,
				Actual:     strconv.Itoa(PresentationActual),
				Percentage: fmt.Sprintf("%d%%", percentage),
			})
		} else if item.Name == "Indiana" {
			performance, err := h.performanceIndianaService.GetDataInByUserAndMonth(userID, month, year)
			if err != nil {
				response := helper.APIResponse("Error fetching Indiana Performance", fiber.StatusInternalServerError, "error", nil)
				return c.Status(fiber.StatusInternalServerError).JSON(response)
			}

			var actual = performance

			target, _ := strconv.Atoi(item.Target) // convertion to int
			percentage := 0
			if target > 0 {
				percentage = (actual * 100) / target
			}

			performances = append(performances, models.UserPerformance{
				Name:       item.Name,
				Target:     item.Target,
				Actual:     strconv.Itoa(actual),
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

func (h *KpiYaeRangeHandler) GetPerformanceUsers(c *fiber.Ctx) error {
	now := time.Now()
	month := uint(now.Month())
	year := uint(now.Year())

	// optional query
	if m := c.Query("month"); m != "" {
		if parsedMonth, err := strconv.Atoi(m); err == nil && parsedMonth >= 1 && parsedMonth <= 12 {
			month = uint(parsedMonth)
		}
	}

	if y := c.Query("year"); y != "" {
		if parsedYear, err := strconv.Atoi(y); err == nil && parsedYear > 0 {
			year = uint(parsedYear)
		}
	}

	// ==========================
	// 1. Ambil KPI master
	// ==========================
	kpiRange, err := h.service.GetKpiYaeRangeByDate(month, year)
	if err != nil {
		return c.Status(500).JSON(
			helper.APIResponse("Failed to fetch KPI", 500, "error", nil),
		)
	}

	var items []Item
	if err := json.Unmarshal([]byte(kpiRange.Target), &items); err != nil {
		return c.Status(500).JSON(
			helper.APIResponse("Invalid KPI config", 500, "error", nil),
		)
	}

	// ==========================
	// 2. Ambil semua user
	// ==========================
	filters := map[string]string{
		"search":      "",
		"order_by":    "id",
		"order":       "DESC",
		"start_date":  "",
		"end_date":    "",
		"user_status": "active",
	}

	userRole := c.Locals("user_role").(string)
	territoryID := c.Locals("territory_id").(int)

	users, _, err := h.userService.GetAllUsers(
		0, // limit
		false,
		1,
		filters,
		[]string{"YAE"},
		false,
		userRole,
		territoryID,
	)
	if err != nil {
		return c.Status(500).JSON(
			helper.APIResponse("Failed to fetch users", 500, "error", nil),
		)
	}

	var result []models.UserKpiPerformance

	// ==========================
	// 3. Loop user → hitung KPI
	// ==========================
	for _, user := range users {

		userID := int(user.ID)
		var performances []models.UserPerformance

		for _, item := range items {

			targetInt, _ := strconv.Atoi(item.Target)
			actual := 0

			switch item.Name {

			case "Visit":
				actual, err = h.visitHistoryService.
					CountVisitHistory(userID, month, year, "")
				if err != nil {
					return c.Status(500).JSON(
						helper.APIResponse("Counting Visit Error", 500, "error", nil),
					)
				}

			case "Presentasi Demo":
				actual, err = h.visitHistoryService.
					CountVisitHistory(userID, month, year, "presentasi_demo")
				if err != nil {
					return c.Status(500).JSON(
						helper.APIResponse("Counting Presentation Error", 500, "error", nil),
					)
				}

			case "Digipos":
				actual, err = h.performanceDigiposService.
					CountPerformanceByUserYaeCode(userID, month, year)
				if err != nil {
					return c.Status(500).JSON(
						helper.APIResponse("Counting Digipos Error", 500, "error", nil),
					)
				}
			case "Indiana":
				actual, err = h.performanceIndianaService.
					GetDataInByUserAndMonth(userID, month, year)
				if err != nil {
					return c.Status(500).JSON(
						helper.APIResponse("Error fetching Indiana Performance", 500, "error", nil),
					)
				}
			}

			percentage := 0
			if targetInt > 0 {
				percentage = (actual * 100) / targetInt
			}

			performances = append(performances, models.UserPerformance{
				Name:       item.Name,
				Target:     item.Target,
				Actual:     strconv.Itoa(actual),
				Percentage: fmt.Sprintf("%d%%", percentage),
			})
		}

		result = append(result, models.UserKpiPerformance{
			UserID:       user.ID,
			Name:         user.Name,
			YaeCode:      user.YaeCode,
			Performances: performances,
		})
	}

	// ==========================
	// 4. Response
	// ==========================
	return c.Status(200).JSON(
		helper.APIResponse(
			"Get Users Performance Successfully",
			200,
			"success",
			map[string]interface{}{
				"users": result,
			},
		),
	)
}
