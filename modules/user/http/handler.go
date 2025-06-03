package http

import (
	"strconv"
	"strings"

	"byu-crm-service/helper"
	authService "byu-crm-service/modules/auth/service"
	"byu-crm-service/modules/user/service"
	"byu-crm-service/modules/user/validation"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service     service.UserService
	authService authService.AuthService
}

func NewUserHandler(service service.UserService, authService authService.AuthService) *UserHandler {
	return &UserHandler{service: service, authService: authService}
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
	intID := c.Locals("user_id").(int)
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
	orderByMostAssignedPic, _ := strconv.ParseBool(c.Query("order_by_most_assigned_pic", "false"))
	rawOnlyRole := c.Query("only_role", "")
	onlyRole := []string{}
	if rawOnlyRole != "" {
		onlyRole = strings.Split(rawOnlyRole, ",")
	}

	userRole := c.Locals("user_role").(string)
	territoryID := c.Locals("territory_id").(int)

	// Call service with filters
	users, total, err := h.service.GetAllUsers(limit, paginate, page, filters, onlyRole, orderByMostAssignedPic, userRole, territoryID)
	if err != nil {
		response := helper.APIResponse(err.Error(), fiber.StatusInternalServerError, "error", nil)
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

func (h *UserHandler) GetUsersResume(c *fiber.Ctx) error {
	// Parse integer and boolean values
	onlyRole := []string{"YAE", "Buddies", "DS"}

	var userRole string
	if queryUserRole := c.Query("user_role"); queryUserRole != "" {
		userRole = queryUserRole
	} else {
		local := c.Locals("user_role")
		if local == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "user_role tidak ditemukan",
			})
		}
		userRoleStr, ok := local.(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "user_role is not a string",
			})
		}
		userRole = userRoleStr
	}
	var territoryID int
	if queryTerritoryID := c.Query("territory_id"); queryTerritoryID != "" {
		var err error
		territoryID, err = strconv.Atoi(queryTerritoryID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid territory_id parameter",
				"error":   err.Error(),
			})
		}
	} else {
		local := c.Locals("territory_id")
		if local == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "territory_id tidak ditemukan",
			})
		}
		territoryIDStr, ok := local.(int)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "territory_id is not a int",
			})
		}
		territoryID = territoryIDStr
	}

	// Call service with filters
	users, err := h.service.GetUsersResume(onlyRole, userRole, territoryID)
	if err != nil {
		response := helper.APIResponse(err.Error(), fiber.StatusInternalServerError, "error", nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Return response
	responseData := map[string]interface{}{
		"users": users,
	}

	response := helper.APIResponse("Get Users Resume Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *UserHandler) UpdateUserProfile(c *fiber.Ctx) error {
	req := new(validation.UpdateProfileRequest)
	if err := c.BodyParser(req); err != nil {
		response := helper.APIResponse("Invalid request", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Request Validation
	errors := validation.ValidateUpdate(req)
	if errors != nil {
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	intID := c.Locals("user_id").(int)
	user, err := h.service.GetUserByID(uint(intID))

	if err != nil {
		response := helper.APIResponse("Failed to fetch user", fiber.StatusInternalServerError, "error", nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	dataUpdate := make(map[string]interface{})
	dataUpdate["name"] = req.Name

	if req.OldPassword != "" || req.NewPassword != "" || req.ConfirmPassword != "" {
		getUser, _ := h.authService.GetUserByKey("email", user.Email)
		// validate old password
		if !h.authService.CheckPassword(req.OldPassword, getUser.Password) {
			errors := map[string]string{
				"old_password": "Password lama tidak sesuai",
			}
			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}

		// Validasi kesamaan password baru dan konfirmasi
		if req.NewPassword != req.ConfirmPassword {
			errors := map[string]string{
				"confirm_password": "Password baru dan konfirmasi password tidak sama",
			}
			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}

		// Kalau valid, masukkan ke dataUpdate
		dataUpdate["password"] = req.NewPassword
	}

	if err := c.BodyParser(user); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	user, err = h.service.UpdateUserProfile(uint(intID), dataUpdate)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to update user",
			"error":   err.Error(),
		})
	}

	responseData := map[string]interface{}{
		"user": user,
	}

	response := helper.APIResponse("Update User Profile Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}
