package http

import (
	permissionService "byu-crm-service/modules/permission/service"
	"byu-crm-service/modules/role/service"
	"byu-crm-service/modules/role/validation"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"byu-crm-service/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type RoleHandler struct {
	roleService       service.RoleService
	permissionService permissionService.PermissionService
	redis             *redis.Client
}

func NewRoleHandler(
	roleService service.RoleService,
	permissionService permissionService.PermissionService,
	redis *redis.Client) *RoleHandler {
	return &RoleHandler{
		roleService:       roleService,
		permissionService: permissionService,
		redis:             redis}
}

func (h *RoleHandler) GetAllRoles(c *fiber.Ctx) error {
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
	cacheKey := fmt.Sprintf("roles:page=%d:limit=%d:paginate=%v", page, limit, paginate)

	for k, v := range filters {
		cacheKey += fmt.Sprintf(":%s=%s", k, v)
	}

	// Coba ambil dari Redis
	cached, err := h.redis.Get(c.Context(), cacheKey).Result()
	if err == nil {
		// Berhasil ambil dari cache
		var cachedData map[string]interface{}
		json.Unmarshal([]byte(cached), &cachedData)

		response := helper.APIResponse("Get Roles (From Cache) Successfully", fiber.StatusOK, "success", cachedData)
		return c.Status(fiber.StatusOK).JSON(response)
	}

	// Call service with filters
	roles, total, err := h.roleService.GetAllRoles(limit, paginate, page, filters)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch roles",
			"error":   err.Error(),
		})
	}

	// Return response
	responseData := map[string]interface{}{
		"roles": roles,
		"total": total,
	}

	// Simpan ke Redis (misal selama 5 menit)
	cacheBytes, _ := json.Marshal(responseData)
	h.redis.Set(c.Context(), cacheKey, cacheBytes, 5*time.Minute)

	response := helper.APIResponse("Get Roles Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *RoleHandler) GetRoleByID(c *fiber.Ctx) error {
	id := c.Params("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid role ID",
			"error":   err.Error(),
		})
	}
	// Generate Redis Cache Key berdasarkan semua filter
	cacheKey := fmt.Sprintf("roleByID:id=%d:", intID)

	// Coba ambil dari Redis
	cached, err := h.redis.Get(c.Context(), cacheKey).Result()
	if err == nil {
		// Berhasil ambil dari cache
		var cachedData map[string]interface{}
		json.Unmarshal([]byte(cached), &cachedData)

		response := helper.APIResponse("Get Role By ID (From Cache) Successfully", fiber.StatusOK, "success", cachedData)
		return c.Status(fiber.StatusOK).JSON(response)
	}

	role, err := h.roleService.GetRoleByID(intID)
	if err != nil {
		response := helper.APIResponse("Failed to fetch role", fiber.StatusInternalServerError, "error", err.Error())
		return c.Status(fiber.StatusOK).JSON(response)
	}

	permissions, err := h.permissionService.GetAllPermissionsByRoleID(intID)
	if err != nil {
		response := helper.APIResponse("Failed to fetch permissions for role", fiber.StatusInternalServerError, "error", err.Error())
		return c.Status(fiber.StatusOK).JSON(response)
	}

	// Return response
	responseData := map[string]interface{}{
		"role":        role,
		"permissions": permissions,
	}

	// Simpan ke Redis (misal selama 5 menit)
	cacheBytes, _ := json.Marshal(responseData)
	h.redis.Set(c.Context(), cacheKey, cacheBytes, 5*time.Minute)

	response := helper.APIResponse("Get Role Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *RoleHandler) CreateRole(c *fiber.Ctx) error {
	req := new(validation.CreateRoleRequest)
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

	existingPermission, _ := h.roleService.GetRoleByName(req.Name)
	if existingPermission != nil {
		errors := map[string]string{
			"name": "Nama role sudah digunakan",
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Konversi permission_id dari []string ke []int
	var permissionIDs []int
	for _, pid := range req.PermissionIDs {
		id, err := strconv.Atoi(pid)
		if err != nil {
			errors := map[string]string{
				"permission_id": "Semua ID permission harus berupa angka",
			}
			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
		permissionIDs = append(permissionIDs, id)
	}

	role, err := h.roleService.CreateRole(&req.Name)
	if err != nil {
		response := helper.APIResponse(err.Error(), fiber.StatusUnauthorized, "error", nil)
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	err = h.permissionService.UpdateRolePermissions(int(role.ID), permissionIDs)
	if err != nil {
		response := helper.APIResponse(err.Error(), fiber.StatusUnauthorized, "error", nil)
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	keys, err := h.redis.Keys(c.Context(), "roles:*").Result()
	if err == nil {
		for _, key := range keys {
			h.redis.Del(c.Context(), key)
		}
	}

	// Response
	response := helper.APIResponse("Roles created successfully", fiber.StatusOK, "success", role)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *RoleHandler) UpdateRole(c *fiber.Ctx) error {
	id := c.Params("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid role ID",
			"error":   err.Error(),
		})
	}
	req := new(validation.UpdateRoleRequest)
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

	currentPermssion, _ := h.roleService.GetRoleByID(intID)
	if currentPermssion == nil {
		errors := map[string]string{
			"name": "Role tidak ditemukan",
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	req.Name = strings.ToLower(strings.TrimSpace(req.Name))

	existingPermission, _ := h.roleService.GetRoleByName(req.Name)

	if existingPermission != nil && currentPermssion.Name != req.Name {
		errors := map[string]string{
			"name": "Nama role sudah digunakan",
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Konversi permission_id dari []string ke []int
	var permissionIDs []int
	for _, pid := range req.PermissionIDs {
		id, err := strconv.Atoi(pid)
		if err != nil {
			errors := map[string]string{
				"permission_id": "Semua ID permission harus berupa angka",
			}
			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
		permissionIDs = append(permissionIDs, id)
	}

	role, err := h.roleService.UpdateRole(&req.Name, intID)
	if err != nil {
		response := helper.APIResponse(err.Error(), fiber.StatusUnauthorized, "error", nil)
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	err = h.permissionService.UpdateRolePermissions(int(role.ID), permissionIDs)
	if err != nil {
		response := helper.APIResponse(err.Error(), fiber.StatusUnauthorized, "error", nil)
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	keys, err := h.redis.Keys(c.Context(), "roles:*").Result()
	if err == nil {
		for _, key := range keys {
			h.redis.Del(c.Context(), key)
		}
	}

	// Hapus cache detail by ID
	h.redis.Del(c.Context(), fmt.Sprintf("roleByID:id=%d:", intID))

	// Response
	response := helper.APIResponse("Permission updated successful", fiber.StatusOK, "success", role)
	return c.Status(fiber.StatusOK).JSON(response)
}
