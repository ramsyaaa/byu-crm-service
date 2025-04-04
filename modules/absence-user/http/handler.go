package http

import (
	"byu-crm-service/modules/absence-user/service"
	"byu-crm-service/modules/absence-user/validation"
	"strconv"

	"byu-crm-service/helper"

	"github.com/gofiber/fiber/v2"
)

type AbsenceUserHandler struct {
	absenceUserService service.AbsenceUserService
}

func NewAbsenceUserHandler(absenceUserService service.AbsenceUserService) *AbsenceUserHandler {
	return &AbsenceUserHandler{absenceUserService: absenceUserService}
}

func (h *AbsenceUserHandler) GetAllAbsenceUsers(c *fiber.Ctx) error {
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
	user_id, _ := strconv.Atoi(c.Query("user_id", "0"))

	// Call service with filters
	absences, total, err := h.absenceUserService.GetAllAbsences(limit, paginate, page, filters, user_id)
	if err != nil {
		response := helper.APIResponse("Failed to fetch absences", fiber.StatusInternalServerError, "error", nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Return response
	responseData := map[string]interface{}{
		"absences": absences,
		"total":    total,
		"page":     page,
	}

	response := helper.APIResponse("Get Absences Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *AbsenceUserHandler) GetAbsenceUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		response := helper.APIResponse("Invalid Absence User ID", 400, "error", nil)
		return c.Status(400).JSON(response)
	}
	AbsenceUser, err := h.absenceUserService.GetAbsenceUserByID(intID)
	if err != nil {
		response := helper.APIResponse("Failed to fetch Absence User", fiber.StatusNotFound, "error", nil)
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	// Return response
	responseData := map[string]interface{}{
		"absence": AbsenceUser,
	}

	response := helper.APIResponse("Get Absence User Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *AbsenceUserHandler) CreateAbsenceUser(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	req := new(validation.CreateAbsenceUserRequest)
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

	data := map[string]any{
		"Visit Account": "App\\Models\\Account",
		"Daily":         nil,
	}

	if !isValidType(req.Type, data) {
		errors := map[string]string{
			"type": "Type tidak valid",
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	subject_type := GetModelValueByKey(data, req.Type)

	var subjectID int
	type_checking := "daily"
	if req.Type == "Visit Account" {
		subjectIDStr := c.FormValue("subject_id")
		subjectID, _ := strconv.Atoi(subjectIDStr)
		type_checking = "monthly"

		if subjectID == 0 {
			errors := map[string]string{
				"subject_id": "subject_id harus diisi",
			}
			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	}

	existingAbsenceUser, message, _ := h.absenceUserService.GetAbsenceUserToday(userID, &req.Type, type_checking)
	if existingAbsenceUser != nil {
		errors := map[string]string{
			"name": message,
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	subjectTypeStr, ok := subject_type.(string)
	if !ok {
		response := helper.APIResponse("Invalid subject type", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	AbsenceUser, err := h.absenceUserService.CreateAbsenceUser(userID, subjectTypeStr, subjectID, &req.Description, &req.Type, &req.Latitude, &req.Longitude)
	if err != nil {
		response := helper.APIResponse(err.Error(), fiber.StatusUnauthorized, "error", AbsenceUser)
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	// Response
	response := helper.APIResponse("Absence user created successful", fiber.StatusOK, "success", AbsenceUser)
	return c.Status(fiber.StatusOK).JSON(response)
}

func GetModelValueByKey(data map[string]any, key string) any {
	if val, ok := data[key]; ok {
		return val
	}
	return nil
}

func isValidType(inputType string, allowed map[string]any) bool {
	_, exists := allowed[inputType]
	return exists
}
