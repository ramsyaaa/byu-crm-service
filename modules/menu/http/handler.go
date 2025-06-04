package http

import (
	"byu-crm-service/constants"
	"byu-crm-service/modules/menu/service"

	"byu-crm-service/helper"

	"github.com/gofiber/fiber/v2"
)

type MenuHandler struct {
	menuService service.MenuService
}

func NewMenuHandler(menuService service.MenuService) *MenuHandler {
	return &MenuHandler{menuService: menuService}
}

func (h *MenuHandler) GetAllMenus(c *fiber.Ctx) error {
	permissionsIface := c.Locals("permissions")

	permissions, ok := permissionsIface.([]string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(helper.APIResponse("Invalid permissions data", fiber.StatusBadRequest, "error", nil))
	}

	// Buat map permission agar cek lebih cepat
	permissionMap := make(map[string]bool)
	for _, perm := range permissions {
		permissionMap[perm] = true
	}

	// Filter menu berdasarkan permission
	var filteredMenus []constants.Menu

	for _, menu := range constants.MenuList {
		var filteredSubMenu []constants.SubMenu

		for _, sub := range menu.SubMenu {
			if permissionMap[sub.Permission] {
				filteredSubMenu = append(filteredSubMenu, sub)
			}
		}

		// Jika setidaknya ada satu sub menu yang cocok, tampilkan menu ini
		if len(filteredSubMenu) > 0 {
			filteredMenus = append(filteredMenus, constants.Menu{
				Name:    menu.Name,
				Icon:    menu.Icon,
				SubMenu: filteredSubMenu,
			})
		}
	}

	responseData := map[string]interface{}{
		"menus": filteredMenus,
	}

	response := helper.APIResponse("Get Menus Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}
