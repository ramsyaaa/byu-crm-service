package http

import (
	"strconv"

	"byu-crm-service/helper"
	"byu-crm-service/modules/product/service"

	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) GetAllProducts(c *fiber.Ctx) error {
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
	userRole := c.Locals("user_role").(string)
	territoryID := c.Locals("territory_id").(int)
	userID := c.Locals("user_id").(int)
	accountID, _ := strconv.Atoi(c.Query("account_id", "1"))

	// Call service with filters
	products, total, err := h.service.GetAllProducts(limit, paginate, page, filters, userRole, territoryID, userID, accountID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch products",
			"error":   err.Error(),
		})
	}

	// Return response
	responseData := map[string]interface{}{
		"products": products,
		"total":    total,
		"page":     page,
	}

	response := helper.APIResponse("Get products Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}
