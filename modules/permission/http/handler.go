package http

import (
	"byu-crm-service/modules/permission/service"
	"byu-crm-service/modules/permission/validation"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"byu-crm-service/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type PermissionHandler struct {
	permissionService service.PermissionService
	redis             *redis.Client
}

func NewPermissionHandler(
	permissionService service.PermissionService,
	redis *redis.Client) *PermissionHandler {
	return &PermissionHandler{permissionService: permissionService,
		redis: redis}
}

func (h *PermissionHandler) GetAllPermissions(c *fiber.Ctx) error {
	// Default query params
	filters := map[string]string{
		"search":     c.Query("search", ""),
		"order_by":   c.Query("order_by", "id"),
		"order":      c.Query("order", "DESC"),
		"start_date": c.Query("start_date", ""),
		"end_date":   c.Query("end_date", ""),
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	paginate, _ := strconv.ParseBool(c.Query("paginate", "true"))
	page, _ := strconv.Atoi(c.Query("page", "1"))

	// Generate Redis Cache Key berdasarkan semua filter
	cacheKey := fmt.Sprintf("permissions:page=%d:limit=%d:paginate=%v", page, limit, paginate)

	for k, v := range filters {
		cacheKey += fmt.Sprintf(":%s=%s", k, v)
	}

	// Coba ambil dari Redis
	cached, err := h.redis.Get(c.Context(), cacheKey).Result()
	if err == nil {
		// Berhasil ambil dari cache
		var cachedData map[string]interface{}
		json.Unmarshal([]byte(cached), &cachedData)

		response := helper.APIResponse("Get Permissions (From Cache) Successfully", fiber.StatusOK, "success", cachedData)
		return c.Status(fiber.StatusOK).JSON(response)
	}

	// Call service with filters
	permissions, total, err := h.permissionService.GetAllPermissions(limit, paginate, page, filters)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch permissions",
			"error":   err.Error(),
		})
	}

	// Return response
	responseData := map[string]interface{}{
		"permissions": permissions,
		"total":       total,
	}

	// Simpan ke Redis (misal selama 5 menit)
	cacheBytes, _ := json.Marshal(responseData)
	h.redis.Set(c.Context(), cacheKey, cacheBytes, 5*time.Minute)

	response := helper.APIResponse("Get Permissions Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *PermissionHandler) GetPermissionByID(c *fiber.Ctx) error {
	id := c.Params("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid permission ID",
			"error":   err.Error(),
		})
	}
	// Generate Redis Cache Key berdasarkan semua filter
	cacheKey := fmt.Sprintf("permisssionByID:id=%d:", intID)

	// Coba ambil dari Redis
	cached, err := h.redis.Get(c.Context(), cacheKey).Result()
	if err == nil {
		// Berhasil ambil dari cache
		var cachedData map[string]interface{}
		json.Unmarshal([]byte(cached), &cachedData)

		response := helper.APIResponse("Get Permission By ID (From Cache) Successfully", fiber.StatusOK, "success", cachedData)
		return c.Status(fiber.StatusOK).JSON(response)
	}

	permission, err := h.permissionService.GetPermissionByID(intID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch permission",
			"error":   err.Error(),
		})
	}

	// Return response
	responseData := map[string]interface{}{
		"permission": permission,
	}

	// Simpan ke Redis (misal selama 5 menit)
	cacheBytes, _ := json.Marshal(responseData)
	h.redis.Set(c.Context(), cacheKey, cacheBytes, 5*time.Minute)

	response := helper.APIResponse("Get Permission Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *PermissionHandler) CreatePermission(c *fiber.Ctx) error {
	req := new(validation.CreatePermissionRequest)
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

	req.Name = strings.ToLower(strings.TrimSpace(req.Name))

	existingPermission, _ := h.permissionService.GetPermissionByName(req.Name)
	if existingPermission != nil {
		errors := map[string]string{
			"name": "Nama permission sudah digunakan",
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	permission, err := h.permissionService.CreatePermission(&req.Name)
	if err != nil {
		response := helper.APIResponse(err.Error(), fiber.StatusUnauthorized, "error", nil)
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	keys, err := h.redis.Keys(c.Context(), "permissions:*").Result()
	if err == nil {
		for _, key := range keys {
			h.redis.Del(c.Context(), key)
		}
	}

	// Response
	response := helper.APIResponse("Permission created successful", fiber.StatusOK, "success", permission)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *PermissionHandler) UpdatePermission(c *fiber.Ctx) error {
	id := c.Params("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid permission ID",
			"error":   err.Error(),
		})
	}
	req := new(validation.UpdatePermissionRequest)
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

	currentPermssion, _ := h.permissionService.GetPermissionByID(intID)
	if currentPermssion == nil {
		errors := map[string]string{
			"name": "Permission tidak ditemukan",
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	req.Name = strings.ToLower(strings.TrimSpace(req.Name))

	existingPermission, _ := h.permissionService.GetPermissionByName(req.Name)

	if existingPermission != nil && currentPermssion.Name != req.Name {
		errors := map[string]string{
			"name": "Nama permission sudah digunakan",
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	permission, err := h.permissionService.UpdatePermission(&req.Name, intID)
	if err != nil {
		response := helper.APIResponse(err.Error(), fiber.StatusUnauthorized, "error", nil)
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	keys, err := h.redis.Keys(c.Context(), "permissions:*").Result()
	if err == nil {
		for _, key := range keys {
			h.redis.Del(c.Context(), key)
		}
	}

	// Hapus cache detail by ID
	h.redis.Del(c.Context(), fmt.Sprintf("permisssionByID:id=%d:", intID))

	// Response
	response := helper.APIResponse("Permission updated successful", fiber.StatusOK, "success", permission)
	return c.Status(fiber.StatusOK).JSON(response)
}
