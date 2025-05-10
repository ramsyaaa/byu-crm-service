package http

import (
	"byu-crm-service/modules/cluster/service"
	"byu-crm-service/modules/cluster/validation"
	"strconv"
	"strings"

	"byu-crm-service/helper"

	"github.com/gofiber/fiber/v2"
)

type ClusterHandler struct {
	clusterService service.ClusterService
}

func NewClusterHandler(clusterService service.ClusterService) *ClusterHandler {
	return &ClusterHandler{clusterService: clusterService}
}

func (h *ClusterHandler) GetAllClusters(c *fiber.Ctx) error {
	// Default query params
	filters := map[string]string{
		"search": c.Query("search", ""),
	}

	// Parse integer and boolean values
	userRole := c.Locals("user_role").(string)
	territoryID := c.Locals("territory_id").(int)

	// Call service with filters
	clusters, total, err := h.clusterService.GetAllClusters(filters, userRole, territoryID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch clusters",
			"error":   err.Error(),
		})
	}

	// Return response
	responseData := map[string]interface{}{
		"clusters": clusters,
		"total":    total,
	}

	response := helper.APIResponse("Get Clusters Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *ClusterHandler) GetClusterByID(c *fiber.Ctx) error {
	id := c.Params("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid cluster ID",
			"error":   err.Error(),
		})
	}
	cluster, err := h.clusterService.GetClusterByID(intID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch cluster",
			"error":   err.Error(),
		})
	}

	// Return response
	responseData := map[string]interface{}{
		"cluster": cluster,
	}

	response := helper.APIResponse("Get Cluster Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *ClusterHandler) CreateCluster(c *fiber.Ctx) error {
	req := new(validation.CreateClusterRequest)
	if err := c.BodyParser(req); err != nil {
		response := helper.APIResponse(err.Error(), fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Request Validation
	errors := validation.ValidateCreate(req)
	if errors != nil {
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	req.Name = strings.ToUpper(strings.TrimSpace(req.Name))

	cluster, err := h.clusterService.CreateCluster(&req.Name, req.BranchID)
	if err != nil {
		response := helper.APIResponse(err.Error(), fiber.StatusUnauthorized, "error", nil)
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	// Response
	response := helper.APIResponse("Cluster created successful", fiber.StatusOK, "success", cluster)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *ClusterHandler) UpdateCluster(c *fiber.Ctx) error {
	id := c.Params("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid cluster ID",
			"error":   err.Error(),
		})
	}
	req := new(validation.UpdateClusterRequest)
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

	currentCluster, _ := h.clusterService.GetClusterByID(intID)
	if currentCluster == nil {
		errors := map[string]string{
			"name": "Cluster tidak ditemukan",
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	req.Name = strings.ToUpper(strings.TrimSpace(req.Name))

	cluster, err := h.clusterService.UpdateCluster(&req.Name, req.BranchID, intID)
	if err != nil {
		response := helper.APIResponse(err.Error(), fiber.StatusUnauthorized, "error", nil)
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	// Response
	response := helper.APIResponse("Cluster updated successful", fiber.StatusOK, "success", cluster)
	return c.Status(fiber.StatusOK).JSON(response)
}
