package http

import (
	"byu-crm-service/modules/absence-user/service"
	"byu-crm-service/modules/absence-user/validation"
	visitHistoryService "byu-crm-service/modules/visit-history/service"
	"strconv"

	"byu-crm-service/helper"

	"github.com/gofiber/fiber/v2"
)

type AbsenceUserHandler struct {
	absenceUserService  service.AbsenceUserService
	visitHistoryService visitHistoryService.VisitHistoryService
}

func NewAbsenceUserHandler(absenceUserService service.AbsenceUserService, visitHistoryService visitHistoryService.VisitHistoryService) *AbsenceUserHandler {
	return &AbsenceUserHandler{absenceUserService: absenceUserService, visitHistoryService: visitHistoryService}
}

func (h *AbsenceUserHandler) GetAllAbsenceUsers(c *fiber.Ctx) error {
	// Default query params
	user_id := c.Locals("user_id").(int)

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
		errors := map[string]string{
			"id": "ID tidak valid",
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
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

	actionType := c.FormValue("action_type")

	if actionType != "Clock In" && actionType != "Clock Out" {
		errors := map[string]string{
			"action_type": "action_type hanya boleh bernilai 'Clock In' atau 'Clock Out'",
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	description := c.FormValue("description")

	if description == "" && actionType == "Clock Out" {
		errors := map[string]string{
			"description": "Deskripsi harus diisi",
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	subject_type := GetModelValueByKey(data, req.Type)
	subjectTypeStr, ok := subject_type.(string)
	if !ok {
		subjectTypeStr = ""
	}

	var subjectID int
	var greeting, survey, presentation *bool
	type_checking := "daily"

	if req.Type == "Visit Account" {
		subjectIDStr := c.FormValue("subject_id")
		parsedSubjectID, _ := strconv.Atoi(subjectIDStr)
		subjectID = parsedSubjectID

		type_checking = "monthly"

		if parsedSubjectID == 0 {
			errors := map[string]string{
				"subject_id": "subject_id harus diisi",
			}
			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
		existingAbsenceUser, message, _ := h.absenceUserService.GetAbsenceUserToday(
			true,
			userID,
			&req.Type,
			type_checking,
			actionType,
			subjectTypeStr,
			parsedSubjectID,
		)
		if existingAbsenceUser != nil {
			errors := map[string]string{
				"message": message,
			}
			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}

		if actionType == "Clock Out" {
			greetingStr := c.FormValue("greeting")
			surveyStr := c.FormValue("survey")
			presentationStr := c.FormValue("presentation")

			errors := make(map[string]string)

			if greetingStr != "" {
				val, _ := strconv.Atoi(greetingStr)
				temp := val != 0
				greeting = &temp
			} else {
				errors["greeting"] = "Salam harus diisi"
			}

			if surveyStr != "" {
				val, _ := strconv.Atoi(surveyStr)
				temp := val != 0
				survey = &temp
			} else {
				errors["survey"] = "Survey harus diisi"
			}

			if presentationStr != "" {
				val, _ := strconv.Atoi(presentationStr)
				temp := val != 0
				presentation = &temp
			} else {
				errors["presentation"] = "Presentasi harus diisi"
			}

			// If there are any errors, return them
			if len(errors) > 0 {
				response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
				return c.Status(fiber.StatusBadRequest).JSON(response)
			}

		}
	} else if req.Type == "Daily" {

		existingAbsenceUser, message, _ := h.absenceUserService.GetAbsenceUserToday(
			true,
			userID,
			&req.Type,
			type_checking,
			actionType,
			"",
			0,
		)
		if existingAbsenceUser != nil {
			errors := map[string]string{
				"message": message,
			}
			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	}

	if actionType == "Clock In" {
		existingAbsenceUser, message, _ := h.absenceUserService.GetAbsenceUserToday(
			false,
			userID,
			&req.Type,
			type_checking,
			actionType,
			"",
			0,
		)

		if existingAbsenceUser != nil {
			errors := map[string]string{
				"message": message,
			}
			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}

		AbsenceUser, err := h.absenceUserService.CreateAbsenceUser(userID, subjectTypeStr, subjectID, &description, &req.Type, &req.Latitude, &req.Longitude)
		if err != nil {
			response := helper.APIResponse(err.Error(), fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
		response := helper.APIResponse("Absence user created successful", fiber.StatusOK, "success", AbsenceUser)
		return c.Status(fiber.StatusOK).JSON(response)
	} else if actionType == "Clock Out" {
		existingAbsenceUser, message, _ := h.absenceUserService.GetAbsenceUserToday(
			false,
			userID,
			&req.Type,
			type_checking,
			actionType,
			"",
			0,
		)

		if existingAbsenceUser == nil {
			errors := map[string]string{
				"message": message,
			}
			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}

		existingAbsenceUser, _, _ = h.absenceUserService.GetAbsenceUserToday(false, userID, &req.Type, type_checking, "Clock In", "", 0)
		if existingAbsenceUser == nil {
			response := helper.APIResponse("No existing absence user found for Clock In", fiber.StatusBadRequest, "error", nil)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
		AbsenceUser, err := h.absenceUserService.UpdateAbsenceUser(int(existingAbsenceUser.ID), userID, subjectTypeStr, subjectID, &description, &req.Type, &req.Latitude, &req.Longitude)
		if err != nil {
			response := helper.APIResponse(err.Error(), fiber.StatusUnauthorized, "error", AbsenceUser)
			return c.Status(fiber.StatusUnauthorized).JSON(response)
		}

		if req.Type == "Visit Account" {
			VisitHistory, err := h.visitHistoryService.CreateVisitHistory(userID, subjectTypeStr, subjectID, int(AbsenceUser.ID), *greeting, *survey, *presentation, &description)
			if err != nil {
				response := helper.APIResponse(err.Error(), fiber.StatusUnauthorized, "error", VisitHistory)
				return c.Status(fiber.StatusUnauthorized).JSON(response)
			}
		}
		response := helper.APIResponse("Absence user created successful", fiber.StatusOK, "success", AbsenceUser)
		return c.Status(fiber.StatusOK).JSON(response)
	}

	// Response
	response := helper.APIResponse("Internal Server Error", fiber.StatusInternalServerError, "success", nil)
	return c.Status(fiber.StatusInternalServerError).JSON(response)
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
