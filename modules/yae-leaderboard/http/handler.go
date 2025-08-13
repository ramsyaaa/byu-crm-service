package http

import (
	"time"

	"byu-crm-service/helper"
	userService "byu-crm-service/modules/user/service"
	"byu-crm-service/modules/yae-leaderboard/service"

	"github.com/gofiber/fiber/v2"
)

type YaeLeaderboardHandler struct {
	service     service.YaeLeaderboardService
	userService userService.UserService
}

func NewYaeLeaderboardHandler(service service.YaeLeaderboardService, userService userService.UserService) *YaeLeaderboardHandler {
	return &YaeLeaderboardHandler{service: service, userService: userService}
}

func (h *YaeLeaderboardHandler) GetAllLeaderboards(c *fiber.Ctx) error {
	// Default query params
	filters := map[string]string{
		"search":      "",
		"order_by":    "id",
		"order":       "DESC",
		"start_date":  "",
		"end_date":    "",
		"user_status": "active",
	}

	// Parse integer and boolean values
	limit := 0
	paginate := false
	page := 1
	orderByMostAssignedPic := false
	onlyRole := []string{"YAE"}

	userRole := c.Locals("user_role").(string)
	territoryID := c.Locals("territory_id").(int)

	var err error

	// Ambil query parameter filter
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	// Parse tanggal jika diisi
	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helper.APIResponse("Format start_date salah (YYYY-MM-DD)", fiber.StatusBadRequest, "error", nil))
		}
	}
	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(helper.APIResponse("Format end_date salah (YYYY-MM-DD)", fiber.StatusBadRequest, "error", nil))
		}
	}

	// Call service with filters
	users, _, err := h.userService.GetAllUsers(limit, paginate, page, filters, onlyRole, orderByMostAssignedPic, userRole, territoryID)
	if err != nil {
		response := helper.APIResponse(err.Error(), fiber.StatusInternalServerError, "error", nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	var userIDs []int
	for _, u := range users {
		userIDs = append(userIDs, int(u.ID))
	}

	// Ambil leaderboard berdasarkan user dan tanggal
	leaderboardData, err := h.service.GetAllLeaderboards(userIDs, startDate, endDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helper.APIResponse("Gagal mengambil leaderboard", fiber.StatusInternalServerError, "error", err.Error()))
	}

	// Return response
	responseData := map[string]interface{}{
		"leaderboard": leaderboardData,
		"page":        page,
	}
	response := helper.APIResponse("Get Leaderboard Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}
