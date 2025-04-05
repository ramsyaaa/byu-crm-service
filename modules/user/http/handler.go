package http

import (
	"fmt"
	"strconv"

	"byu-crm-service/helper"
	"byu-crm-service/modules/user/service"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid user ID",
			"error":   err.Error(),
		})
	}
	user, err := h.service.GetUserByID(uint(intID))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch user",
			"error":   err.Error(),
		})
	}

	// Return response
	responseData := map[string]interface{}{
		"user": user,
	}

	response := helper.APIResponse("Get User Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *UserHandler) GetUserProfile(c *fiber.Ctx) error {
	fmt.Println("GetUserProfile called")
	intID := c.Locals("user_id").(int)
	fmt.Println("User fetched successfully:")
	user, err := h.service.GetUserByID(uint(intID))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch user",
			"error":   err.Error(),
		})
	}

	fmt.Println("User fetched successfully:", user)

	// Return response
	responseData := map[string]interface{}{
		"user": user,
	}

	response := helper.APIResponse("Get User Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	// Default query params
	filters := map[string]string{
		"search":     c.Query("search", ""),
		"order_by":   c.Query("order_by", "id"),
		"order":      c.Query("order", "DESC"),
		"start_date": c.Query("start_date", ""),
		"end_date":   c.Query("end_date", ""),
	}

	// Parse integer and boolean values
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	paginate, _ := strconv.ParseBool(c.Query("paginate", "true"))
	page, _ := strconv.Atoi(c.Query("page", "1"))

	// Call service with filters
	users, total, err := h.service.GetAllUsers(limit, paginate, page, filters)
	if err != nil {
		response := helper.APIResponse(err.Error(), fiber.StatusInternalServerError, "success", nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Return response
	responseData := map[string]interface{}{
		"users": users,
		"total": total,
		"page":  page,
	}

	response := helper.APIResponse("Get Users Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}
