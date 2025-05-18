package http

import (
	"byu-crm-service/modules/absence-user/service"
	"byu-crm-service/modules/absence-user/validation"
	accountService "byu-crm-service/modules/account/service"
	kpiYaeRange "byu-crm-service/modules/kpi-yae-range/service"
	visitChecklistService "byu-crm-service/modules/visit-checklist/service"
	visitHistoryService "byu-crm-service/modules/visit-history/service"
	"fmt"
	"strconv"
	"strings"

	"byu-crm-service/helper"

	"github.com/gofiber/fiber/v2"
)

type AbsenceUserHandler struct {
	absenceUserService    service.AbsenceUserService
	visitHistoryService   visitHistoryService.VisitHistoryService
	accountService        accountService.AccountService
	KpiYaeRangeService    kpiYaeRange.KpiYaeRangeService
	visitChecklistService visitChecklistService.VisitChecklistService
}

func NewAbsenceUserHandler(
	absenceUserService service.AbsenceUserService,
	visitHistoryService visitHistoryService.VisitHistoryService,
	accountService accountService.AccountService,
	kpiYaeRange kpiYaeRange.KpiYaeRangeService,
	visitChecklistService visitChecklistService.VisitChecklistService) *AbsenceUserHandler {
	return &AbsenceUserHandler{
		absenceUserService:    absenceUserService,
		visitHistoryService:   visitHistoryService,
		accountService:        accountService,
		KpiYaeRangeService:    kpiYaeRange,
		visitChecklistService: visitChecklistService,
	}
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
	month, _ := strconv.Atoi(c.Query("month", "0"))
	year, _ := strconv.Atoi(c.Query("year", "0"))
	absence_type := c.Query("type", "")

	// Call service with filters
	absences, total, err := h.absenceUserService.GetAllAbsences(limit, paginate, page, filters, user_id, month, year, absence_type)
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
	territoryID := c.Locals("territory_id").(int)
	userRole := c.Locals("user_role").(string)
	helper.LogError(c, fmt.Sprintf("Failed to create absence: %v", "before validate"))
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

	helper.LogError(c, fmt.Sprintf("Failed to create absence: %v", "after validate"))

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
	kpiYae := make(map[string]*int)
	detailVisit := make(map[string]string)

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
	type_checking := "daily"
	helper.LogError(c, fmt.Sprintf("Failed to create absence: %v", "tes1"))

	if req.Type == "Visit Account" {
		helper.LogError(c, fmt.Sprintf("Failed to create absence: %v", "tes2"))
		subjectIDStr := c.FormValue("subject_id")
		parsedSubjectID, _ := strconv.Atoi(subjectIDStr)
		subjectID = parsedSubjectID

		existingAbsenceUser, _ := h.absenceUserService.AlreadyAbsenceInSameDay(
			userID,
			&req.Type,
			type_checking,
			actionType,
			subjectTypeStr,
			subjectID,
		)
		if actionType == "Clock In" {
			helper.LogError(c, fmt.Sprintf("Failed to create absence: %v", "tes3"))
			if existingAbsenceUser != nil {
				errors := map[string]string{
					"message": "User Already absence today",
				}
				response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
				return c.Status(fiber.StatusBadRequest).JSON(response)
			}

		}

		type_checking = "monthly"

		if parsedSubjectID == 0 {
			errors := map[string]string{
				"subject_id": "subject_id harus diisi",
			}
			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}

		if actionType == "Clock In" {
			helper.LogError(c, fmt.Sprintf("Failed to create absence: %v", "tes4"))
			getAccount, err := h.accountService.FindByAccountID(uint(parsedSubjectID), userRole, uint(territoryID), uint(userID))
			if err != nil {
				response := helper.APIResponse("Failed to fetch account", fiber.StatusInternalServerError, "error", nil)
				return c.Status(fiber.StatusInternalServerError).JSON(response)
			}

			if userRole != "Super-Admin" {
				if getAccount.Latitude != nil && getAccount.Longitude != nil &&
					*getAccount.Latitude != "" && *getAccount.Longitude != "" {
					latitude, err := strconv.ParseFloat(req.Latitude, 64)
					if err != nil {
						response := helper.APIResponse("Invalid latitude value", fiber.StatusBadRequest, "error", nil)
						return c.Status(fiber.StatusBadRequest).JSON(response)
					}
					longitude, err := strconv.ParseFloat(req.Longitude, 64)
					if err != nil {
						response := helper.APIResponse("Invalid longitude value", fiber.StatusBadRequest, "error", nil)
						return c.Status(fiber.StatusBadRequest).JSON(response)
					}
					accountLatitude, err := strconv.ParseFloat(*getAccount.Latitude, 64)
					if err != nil {
						response := helper.APIResponse("Invalid latitude value in account", fiber.StatusBadRequest, "error", nil)
						return c.Status(fiber.StatusBadRequest).JSON(response)
					}
					accountLongitude, err := strconv.ParseFloat(*getAccount.Longitude, 64)
					if err != nil {
						response := helper.APIResponse("Invalid longitude value in account", fiber.StatusBadRequest, "error", nil)
						return c.Status(fiber.StatusBadRequest).JSON(response)
					}
					inRadius := helper.IsWithinRadius(100, latitude, longitude, accountLatitude, accountLongitude)
					if !inRadius {
						errors := map[string]string{
							"radius": "Anda tidak berada dalam radius 100 meter dari lokasi account / data longitude dan latitude tidak valid",
						}
						response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
						return c.Status(fiber.StatusBadRequest).JSON(response)
					}
				} else {
					requestBody := map[string]interface{}{
						"longitude": req.Longitude,
						"latitude":  req.Latitude,
					}
					_, err := h.accountService.UpdateAccount(requestBody, parsedSubjectID, userRole, territoryID, userID)

					if err != nil {
						response := helper.APIResponse(err.Error(), fiber.StatusInternalServerError, "error", nil)
						return c.Status(fiber.StatusInternalServerError).JSON(response)
					}
				}
			}
		} else if actionType == "Clock Out" {
			helper.LogError(c, fmt.Sprintf("Failed to create absence: %v", "tes5"))
			getVisitList, err := h.visitChecklistService.GetAllVisitChecklist()

			if err != nil {
				response := helper.APIResponse("Failed to fetch visit list", fiber.StatusInternalServerError, "error", nil)
				return c.Status(fiber.StatusInternalServerError).JSON(response)
			}

			errors := make(map[string]string)

			for _, item := range getVisitList {
				formKey := item.Key
				// nameKey := item.Name
				valueBytes := c.Context().FormValue(formKey)
				valueStr := string(valueBytes)

				if valueStr == "" {
					errors[formKey] = fmt.Sprintf("%s harus diisi", item.Name)
					continue
				}

				if valueStr != "1" && valueStr != "0" {
					errors[formKey] = fmt.Sprintf("%s hanya boleh bernilai 1 atau 0", item.Name)
					continue
				}

				parsedValue, err := strconv.Atoi(valueStr)
				if err != nil {
					errors[formKey] = fmt.Sprintf("%s harus berupa angka", item.Name)
					continue
				}

				if formKey == "presentasi_demo" {
					if valueStr == "1" {
						demoDocumentation := c.FormValue("demo_documentation")
						if demoDocumentation == "" {
							errors := map[string]string{
								"demo_documentation": "Dokumentasi demo harus diisi",
							}
							response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
							return c.Status(fiber.StatusBadRequest).JSON(response)
						}

						// Simpan gambar base64
						filePath, err := helper.SaveValidatedBase64File(demoDocumentation, "public/uploads/demo")
						if err != nil {
							errors := map[string]string{
								"demo_documentation": err.Error(),
							}
							response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
							return c.Status(fiber.StatusBadRequest).JSON(response)
						}

						// Simpan relative path
						detailVisit[formKey] = strings.TrimPrefix(filePath, "public/")
					} else if valueStr == "0" {
						demoReason := c.FormValue("demo_reason")
						if demoReason == "" {
							errors := map[string]string{
								"demo_reason": "Alasan tidak presentasi / demo harus diisi",
							}
							response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
							return c.Status(fiber.StatusBadRequest).JSON(response)
						}

						detailVisit[formKey+"_reason"] = demoReason
					}
				}

				if formKey == "dealing_sekolah" {
					if valueStr == "1" {
						bakFile := c.FormValue("bak_file")
						if bakFile == "" {
							errors := map[string]string{
								"bak_file": "File BAK harus diisi",
							}
							response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
							return c.Status(fiber.StatusBadRequest).JSON(response)
						}

						// Simpan file BAK base64
						filePath, err := helper.SaveValidatedBase64File(bakFile, "public/uploads/bak")
						if err != nil {
							errors := map[string]string{
								"bak_file": err.Error(),
							}
							response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
							return c.Status(fiber.StatusBadRequest).JSON(response)
						}

						// Simpan relative path
						detailVisit[formKey] = strings.TrimPrefix(filePath, "public/")
					} else if valueStr == "0" {
						dealingReason := c.FormValue("dealing_reason")
						if dealingReason == "" {
							errors := map[string]string{
								"dealing_reason": "Alasan tidak dealing harus diisi",
							}
							response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
							return c.Status(fiber.StatusBadRequest).JSON(response)
						}

						detailVisit[formKey+"_reason"] = dealingReason
					}
				}

				fmt.Println(detailVisit)

				kpiYae[formKey] = &parsedValue
			}

			if len(errors) > 0 {
				response := helper.APIResponse("Validasi gagal", fiber.StatusBadRequest, "error", errors)
				return c.Status(fiber.StatusBadRequest).JSON(response)
			}
		}
	} else if req.Type == "Daily" {

		existingAbsenceUser, _ := h.absenceUserService.AlreadyAbsenceInSameDay(
			userID,
			&req.Type,
			type_checking,
			actionType,
			"",
			0,
		)
		if actionType == "Clock In" {
			if existingAbsenceUser != nil {
				errors := map[string]string{
					"message": "User Already absence today",
				}
				response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
				return c.Status(fiber.StatusBadRequest).JSON(response)
			}

		}

	}

	if actionType == "Clock In" {
		helper.LogError(c, fmt.Sprintf("Failed to create absence: %v", "tes6"))
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
		AbsenceUser, err := h.absenceUserService.UpdateAbsenceUser(int(existingAbsenceUser.ID), userID, subjectTypeStr, subjectID, &description, &req.Type)
		if err != nil {
			response := helper.APIResponse(err.Error(), fiber.StatusUnauthorized, "error", AbsenceUser)
			return c.Status(fiber.StatusUnauthorized).JSON(response)
		}

		if req.Type == "Visit Account" {
			VisitHistory, err := h.visitHistoryService.CreateVisitHistory(userID, subjectTypeStr, subjectID, int(AbsenceUser.ID), kpiYae, &description, detailVisit)
			if err != nil {
				response := helper.APIResponse(err.Error(), fiber.StatusUnauthorized, "error", VisitHistory)
				return c.Status(fiber.StatusUnauthorized).JSON(response)
			}
		}
		helper.LogError(c, fmt.Sprintf("Failed to create absence: %v", "tes7"))

		response := helper.APIResponse("Absence user created successful", fiber.StatusOK, "success", AbsenceUser)
		return c.Status(fiber.StatusOK).JSON(response)
	}

	// Response
	response := helper.APIResponse("Internal Server Error", fiber.StatusInternalServerError, "success", nil)
	return c.Status(fiber.StatusInternalServerError).JSON(response)
}

func (h *AbsenceUserHandler) GetAbsenceActive(c *fiber.Ctx) error {
	// Default query params
	user_id := c.Locals("user_id").(int)

	type_absence := c.Query("type", "")

	// Call service with filters
	absences, err := h.absenceUserService.GetAbsenceActive(user_id, type_absence)
	if err != nil {
		response := helper.APIResponse("Failed to fetch absences", fiber.StatusInternalServerError, "error", nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Return response
	responseData := map[string]interface{}{
		"absences": absences,
	}

	response := helper.APIResponse("Get Absences Successfully", fiber.StatusOK, "success", responseData)
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
