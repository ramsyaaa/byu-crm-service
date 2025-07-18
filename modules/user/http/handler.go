package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"byu-crm-service/helper"
	"byu-crm-service/models"
	accountResponse "byu-crm-service/modules/account/response"
	accountService "byu-crm-service/modules/account/service"
	authService "byu-crm-service/modules/auth/service"
	roleService "byu-crm-service/modules/role/service"
	"byu-crm-service/modules/user/response"
	"byu-crm-service/modules/user/service"
	"byu-crm-service/modules/user/validation"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type UserHandler struct {
	service        service.UserService
	authService    authService.AuthService
	accountService accountService.AccountService
	roleService    roleService.RoleService
	redis          *redis.Client
}

func NewUserHandler(service service.UserService, authService authService.AuthService, accountService accountService.AccountService, roleService roleService.RoleService, redis *redis.Client) *UserHandler {
	return &UserHandler{service: service, authService: authService, accountService: accountService, roleService: roleService, redis: redis}
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

	// Generate Redis Cache Key berdasarkan semua filter
	cacheKey := fmt.Sprintf("userbyID:user=%d", intID)

	// Coba ambil dari Redis
	cached, err := h.redis.Get(c.Context(), cacheKey).Result()
	if err == nil {
		// Berhasil ambil dari cache
		var cachedData map[string]interface{}
		json.Unmarshal([]byte(cached), &cachedData)

		response := helper.APIResponse("Get User By ID (From Cache) Successfully", fiber.StatusOK, "success", cachedData)
		return c.Status(fiber.StatusOK).JSON(response)
	}

	user, err := h.service.GetUserByID(uint(intID))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch user",
			"error":   err.Error(),
		})
	}

	accountsData := []accountResponse.AccountResponse{}
	if user != nil {
		filters := map[string]string{
			"search":           c.Query("search", ""),
			"order_by":         c.Query("order_by", "id"),
			"order":            c.Query("order", "DESC"),
			"start_date":       c.Query("start_date", ""),
			"end_date":         c.Query("end_date", ""),
			"account_category": c.Query("account_category", ""),
			"account_type":     c.Query("account_type", ""),
			"only_skulid":      c.Query("only_skulid", "0"),
			"is_priority":      c.Query("is_priority", "0"),
		}

		// Parse integer and boolean values
		limit := 0
		paginate := false
		page := 1
		userRole := c.Locals("user_role").(string)
		territoryID := c.Locals("territory_id").(int)

		if user.UserType == "Administrator" {
			userRole = "Super-Admin"
		} else if user.UserType == "HQ" {
			userRole = "HQ"
		} else if user.UserType == "AREA" {
			userRole = "Area"
		} else if user.UserType == "REGIONAL" {
			userRole = "Regional"
		} else if user.UserType == "BRANCH" {
			userRole = "Branch"
		}

		territoryID = int(user.TerritoryID)

		onlyUserPic := true
		accounts, _, err := h.accountService.GetAllAccounts(limit, paginate, page, filters, userRole, territoryID, intID, onlyUserPic, false)
		if err != nil {
			response := helper.APIResponse("Failed to fetch account PICs", fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
		accountsData = accounts
	} else {
		responseData := map[string]interface{}{
			"user":     user,
			"accounts": accountsData,
		}

		response := helper.APIResponse("User not found", fiber.StatusNotFound, "error", responseData)
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	// Return response
	responseData := map[string]interface{}{
		"user":     user,
		"accounts": accountsData,
	}

	// Simpan ke Redis (misal selama 5 menit)
	cacheBytes, _ := json.Marshal(responseData)
	h.redis.Set(c.Context(), cacheKey, cacheBytes, 5*time.Minute)

	response := helper.APIResponse("Get User Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	// Add a timeout context to prevent long-running operations
	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	// Use a recovery function to catch any panics
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic in Create User: %v", r)
			response := helper.APIResponse("Internal server error", fiber.StatusInternalServerError, "error", r)
			c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}()

	// Parse request body with error handling
	req := new(validation.ValidateRequest)

	if err := c.BodyParser(req); err != nil {
		// Check for specific EOF error
		if err.Error() == "unexpected EOF" {
			response := helper.APIResponse("Invalid request: Unexpected end of JSON input", fiber.StatusBadRequest, "error", nil)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}

		response := helper.APIResponse("Invalid request format: "+err.Error(), fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	// Request Validation with context
	select {
	case <-ctx.Done():
		response := helper.APIResponse("Request timeout during validation", fiber.StatusRequestTimeout, "error", nil)
		return c.Status(fiber.StatusRequestTimeout).JSON(response)
	default:
		errors := validation.ValidateCreate(req)
		if errors != nil {
			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	}

	errors := validation.AdditionalValidate(req, 0)
	if errors != nil {
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Create Account with context and error handling
	reqMap := make(map[string]interface{})

	// Marshal request to JSON with timeout
	var reqBytes []byte
	var marshalErr error

	select {
	case <-ctx.Done():
		response := helper.APIResponse("Request timeout during marshaling", fiber.StatusRequestTimeout, "error", nil)
		return c.Status(fiber.StatusRequestTimeout).JSON(response)
	default:
		reqBytes, marshalErr = json.Marshal(req)
		if marshalErr != nil {
			log.Printf(fmt.Sprintf("Failed to marshal request: %v", marshalErr))
			response := helper.APIResponse("Failed to process request data", fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	// Unmarshal JSON to map with timeout
	var unmarshalErr error

	select {
	case <-ctx.Done():
		response := helper.APIResponse("Request timeout during unmarshaling", fiber.StatusRequestTimeout, "error", nil)
		return c.Status(fiber.StatusRequestTimeout).JSON(response)
	default:
		unmarshalErr = json.Unmarshal(reqBytes, &reqMap)
		if unmarshalErr != nil {
			log.Printf(fmt.Sprintf("Failed to unmarshal request: %v", unmarshalErr))
			response := helper.APIResponse("Failed to process request data", fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	// Call service with timeout
	var user *models.User
	var serviceErr error

	select {
	case <-ctx.Done():
		response := helper.APIResponse("Request timeout during user creation", fiber.StatusRequestTimeout, "error", nil)
		return c.Status(fiber.StatusRequestTimeout).JSON(response)
	default:
		user, serviceErr = h.service.CreateUser(reqMap)
		if serviceErr != nil {
			log.Printf(fmt.Sprintf("Failed to create user: %v", serviceErr))
			response := helper.APIResponse("Failed to create user: "+serviceErr.Error(), fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	var roleID int
	switch v := reqMap["role_id"].(type) {
	case float64:
		roleID = int(v)
	case int:
		roleID = v
	case string:
		if parsed, err := strconv.Atoi(v); err == nil {
			roleID = parsed
		}
	}
	_ = h.roleService.AssignModelHasRole("App\\Models\\User", int(user.ID), roleID)

	// Return success response
	response := helper.APIResponse("Create User Succsesfully", fiber.StatusOK, "success", user)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	// Add a timeout context to prevent long-running operations
	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	// Use a recovery function to catch any panics
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic in Update User: %v", r)
			response := helper.APIResponse("Internal server error", fiber.StatusInternalServerError, "error", r)
			c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}()

	// Get and validate user ID
	userIDStr := c.Params("id")
	if userIDStr == "" {
		response := helper.APIResponse("User ID is required", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		response := helper.APIResponse("Invalid User ID", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Parse request body with error handling
	req := new(validation.ValidateRequest)

	if err := c.BodyParser(req); err != nil {
		// Check for specific EOF error
		if err.Error() == "unexpected EOF" {
			response := helper.APIResponse("Invalid request: Unexpected end of JSON input", fiber.StatusBadRequest, "error", nil)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}

		response := helper.APIResponse("Invalid request format: "+err.Error(), fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	// Request Validation with context
	select {
	case <-ctx.Done():
		response := helper.APIResponse("Request timeout during validation", fiber.StatusRequestTimeout, "error", nil)
		return c.Status(fiber.StatusRequestTimeout).JSON(response)
	default:
		errors := validation.ValidateCreate(req)
		if errors != nil {
			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	}

	errors := validation.AdditionalValidate(req, userID)
	if errors != nil {
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Create User with context and error handling
	reqMap := make(map[string]interface{})

	// Marshal request to JSON with timeout
	var reqBytes []byte
	var marshalErr error

	select {
	case <-ctx.Done():
		response := helper.APIResponse("Request timeout during marshaling", fiber.StatusRequestTimeout, "error", nil)
		return c.Status(fiber.StatusRequestTimeout).JSON(response)
	default:
		reqBytes, marshalErr = json.Marshal(req)
		if marshalErr != nil {
			log.Printf(fmt.Sprintf("Failed to marshal request: %v", marshalErr))
			response := helper.APIResponse("Failed to process request data", fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	// Unmarshal JSON to map with timeout
	var unmarshalErr error

	select {
	case <-ctx.Done():
		response := helper.APIResponse("Request timeout during unmarshaling", fiber.StatusRequestTimeout, "error", nil)
		return c.Status(fiber.StatusRequestTimeout).JSON(response)
	default:
		unmarshalErr = json.Unmarshal(reqBytes, &reqMap)
		if unmarshalErr != nil {
			log.Printf(fmt.Sprintf("Failed to unmarshal request: %v", unmarshalErr))
			response := helper.APIResponse("Failed to process request data", fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	// Call service with timeout
	var user *response.UserResponse
	var serviceErr error

	select {
	case <-ctx.Done():
		response := helper.APIResponse("Request timeout during user creation", fiber.StatusRequestTimeout, "error", nil)
		return c.Status(fiber.StatusRequestTimeout).JSON(response)
	default:
		user, serviceErr = h.service.UpdateUser(reqMap, userID)
		if serviceErr != nil {
			log.Printf(fmt.Sprintf("Failed to update user: %v", serviceErr))
			response := helper.APIResponse("Failed to update user: "+serviceErr.Error(), fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	var roleID int
	switch v := reqMap["role_id"].(type) {
	case float64:
		roleID = int(v)
	case int:
		roleID = v
	case string:
		if parsed, err := strconv.Atoi(v); err == nil {
			roleID = parsed
		}
	}
	_ = h.roleService.AssignModelHasRole("App\\Models\\User", int(user.ID), roleID)

	// Return success response
	response := helper.APIResponse("Update User Succsesfully", fiber.StatusOK, "success", user)
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
		"search":      c.Query("search", ""),
		"order_by":    c.Query("order_by", "id"),
		"order":       c.Query("order", "DESC"),
		"start_date":  c.Query("start_date", ""),
		"end_date":    c.Query("end_date", ""),
		"user_status": c.Query("user_status", "active"),
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
	errors := validation.ValidateProfile(req)
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
	dataUpdate["msisdn"] = NormalizeMsisdn(req.Msisdn)

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

	userByMsisdn, err := h.service.GetUserByMsisdn(dataUpdate["msisdn"].(string))
	if err != nil {
		response := helper.APIResponse("Failed to fetch user by msisdn", fiber.StatusInternalServerError, "error", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	if userByMsisdn != nil {
		if userByMsisdn.ID != user.ID {
			if userByMsisdn.Msisdn == dataUpdate["msisdn"] {
				errors := map[string]string{
					"msisdn": "MSISDN sudah digunakan oleh pengguna lain",
				}
				response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
				return c.Status(fiber.StatusBadRequest).JSON(response)
			}
		}
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

func NormalizeMsisdn(msisdn string) string {
	msisdn = strings.TrimSpace(msisdn)

	if strings.HasPrefix(msisdn, "+62") {
		return msisdn
	} else if strings.HasPrefix(msisdn, "62") {
		return "+" + msisdn
	} else if strings.HasPrefix(msisdn, "0") {
		return "+62" + msisdn[1:]
	} else if strings.HasPrefix(msisdn, "8") {
		return "+62" + msisdn
	}

	return msisdn
}
