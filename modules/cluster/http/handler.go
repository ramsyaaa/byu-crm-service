package http

import (
	"strconv"

	"byu-crm-service/modules/cluster/service"

	"github.com/gofiber/fiber/v2"
)

type ClusterHandler struct {
	service service.ClusterService
}

func NewClusterHandler(service service.ClusterService) *ClusterHandler {
	return &ClusterHandler{service: service}
}

func (h *ClusterHandler) GetClusterByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID parameter"})
	}

	cluster, err := h.service.GetClusterByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if cluster == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Cluster not found"})
	}

	return c.JSON(cluster)
}

func (h *ClusterHandler) GetClusterByName(c *fiber.Ctx) error {
	name := c.Query("name")
	if name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cluster name is required"})
	}

	cluster, err := h.service.GetClusterByName(name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if cluster == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Cluster not found"})
	}

	return c.JSON(cluster)
}
