package http

import (
	"byu-crm-service/modules/absence-user/service"
	"byu-crm-service/modules/absence-user/validation"
	accountService "byu-crm-service/modules/account/service"
	kpiYaeRange "byu-crm-service/modules/kpi-yae-range/service"
	notificationService "byu-crm-service/modules/notification/service"
	smsSenderService "byu-crm-service/modules/sms-sender/service"
	userService "byu-crm-service/modules/user/service"
	visitChecklistService "byu-crm-service/modules/visit-checklist/service"
	visitHistoryService "byu-crm-service/modules/visit-history/service"
	"fmt"
	"strconv"
	"time"

	"byu-crm-service/helper"

	"github.com/gofiber/fiber/v2"
)

type AbsenceUserHandler struct {
	absenceUserService    service.AbsenceUserService
	visitHistoryService   visitHistoryService.VisitHistoryService
	accountService        accountService.AccountService
	KpiYaeRangeService    kpiYaeRange.KpiYaeRangeService
	visitChecklistService visitChecklistService.VisitChecklistService
	notificationService   notificationService.NotificationService
	smsSenderService      smsSenderService.SmsSenderService
	userService           userService.UserService
}

func NewAbsenceUserHandler(
	absenceUserService service.AbsenceUserService,
	visitHistoryService visitHistoryService.VisitHistoryService,
	accountService accountService.AccountService,
	kpiYaeRange kpiYaeRange.KpiYaeRangeService,
	visitChecklistService visitChecklistService.VisitChecklistService,
	notificationService notificationService.NotificationService,
	smsSender smsSenderService.SmsSenderService,
	userService userService.UserService) *AbsenceUserHandler {
	return &AbsenceUserHandler{
		absenceUserService:    absenceUserService,
		visitHistoryService:   visitHistoryService,
		accountService:        accountService,
		KpiYaeRangeService:    kpiYaeRange,
		visitChecklistService: visitChecklistService,
		notificationService:   notificationService,
		smsSenderService:      smsSender,
		userService:           userService,
	}
}

func (h *AbsenceUserHandler) GetAllAbsenceUsers(c *fiber.Ctx) error {
	// Default query params
	user_id := c.Locals("user_id").(int)

	paramUserID := c.Query("user_id")
	if paramUserID != "" {
		if parsedID, err := strconv.Atoi(paramUserID); err == nil {
			user_id = parsedID
		}
	}

	filters := map[string]string{
		"search":     c.Query("search", ""),
		"order_by":   c.Query("order_by", "id"),
		"order":      c.Query("order", "DESC"),
		"start_date": c.Query("start_date", ""),
		"end_date":   c.Query("end_date", ""),
		"status":     c.Query("status", "1"),
		"all_user":   c.Query("all_user", "0"),
	}

	// Parse integer and boolean values
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	paginate, _ := strconv.ParseBool(c.Query("paginate", "true"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	month, _ := strconv.Atoi(c.Query("month", "0"))
	year, _ := strconv.Atoi(c.Query("year", "0"))
	absence_type := c.Query("type", "")
	userRole := c.Locals("user_role").(string)
	territory_id := c.Locals("territory_id").(int)

	userIDs := []int{}

	if filters["status"] == "0" {
		getUsers, _, err := h.userService.GetAllUsers(0, false, 1, map[string]string{"search": filters["search"]}, []string{}, false, userRole, territory_id)
		if err != nil {
			response := helper.APIResponse("Failed to fetch users", fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}

		for _, user := range getUsers {
			if user.ID == 356 {
				fmt.Println("Skipping user with ID 356")
			}
			userIDs = append(userIDs, int(user.ID))
		}
	}

	// Call service with filters
	absences, total, err := h.absenceUserService.GetAllAbsences(limit, paginate, page, filters, user_id, month, year, absence_type, userRole, territory_id, userIDs)
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
	userID := c.Locals("user_id").(int)
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

	_ = h.notificationService.MarkNotificationAsReadBySubjectID("App\\Models\\AbsenceUser", uint(AbsenceUser.ID), userID)

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
	notifyUser := false

	var successCode int = fiber.StatusOK
	var successMessage string = "Absence user created successfully"

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

	if req.Type == "Visit Account" {
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
			getAccount, err := h.accountService.FindByAccountID(uint(parsedSubjectID), userRole, uint(territoryID), uint(userID))
			if err != nil {
				response := helper.APIResponse("Failed to fetch account", fiber.StatusInternalServerError, "error", nil)
				return c.Status(fiber.StatusInternalServerError).JSON(response)
			}
			typeClockIn := c.FormValue("type_clock_in")

			if typeClockIn == "image" {
				evidence_image := c.FormValue("evidence_image")

				if evidence_image == "" {
					errors := map[string]string{
						"evidence_image": "Gambar harus diisi",
					}
					response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
					return c.Status(fiber.StatusBadRequest).JSON(response)
				}

				decoded, _, err := helper.DecodeBase64Image(evidence_image)
				if err != nil {
					errors := map[string]string{
						"evidence_image": "Gambar tidak valid: " + err.Error(),
					}
					response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
					return c.Status(fiber.StatusBadRequest).JSON(response)
				}

				// Validasi ukuran maksimum 5MB
				if len(decoded) > 5*1024*1024 {
					errors := map[string]string{
						"evidence_image": "Ukuran file maksimum adalah 5MB",
					}
					response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
					return c.Status(fiber.StatusBadRequest).JSON(response)
				}

				// successCode = fiber.StatusCreated
				// successMessage = "Absence user created successfully with image evidence"
				successCode = fiber.StatusOK
				successMessage = "Absence user created successfully"
			} else {
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
						inRadius := helper.IsWithinRadius(500, latitude, longitude, accountLatitude, accountLongitude)
						if !inRadius {
							errors := map[string]string{
								"radius": "Anda tidak berada dalam radius 500 meter dari lokasi account / data longitude dan latitude tidak valid",
							}
							response := helper.APIResponse("Validation error", fiber.StatusUnprocessableEntity, "error", errors)
							return c.Status(fiber.StatusUnprocessableEntity).JSON(response)
						}
					} else {
						// successMessage = "Absence user created successfully, but location data is not available for this account"
						// successCode = fiber.StatusCreated
						successCode = fiber.StatusOK
						successMessage = "Absence user created successfully"
					}
				}
			}
		} else if actionType == "Clock Out" {
			getAccount, err := h.accountService.FindByAccountID(uint(parsedSubjectID), userRole, uint(territoryID), uint(userID))
			if err != nil {
				response := helper.APIResponse("Failed to fetch account", fiber.StatusInternalServerError, "error", nil)
				return c.Status(fiber.StatusInternalServerError).JSON(response)
			}

			getVisitList, err := h.visitChecklistService.GetAllVisitChecklist()

			if err != nil {
				response := helper.APIResponse("Failed to fetch visit list", fiber.StatusInternalServerError, "error", nil)
				return c.Status(fiber.StatusInternalServerError).JSON(response)
			}

			errors := make(map[string]string)
			var statusSkulid string

			if getAccount.IsSkulid == nil || *getAccount.IsSkulid != 1 {
				statusSkulid = c.FormValue("status_skulid")
				if statusSkulid == "" {
					errors["status_skulid"] = "status SkulID harus diisi"
				} else if statusSkulid != "Ya" && statusSkulid != "Belum" && statusSkulid != "Sedang Tahap Penawaran" {
					errors["status_skulid"] = "Harus memilih status hanya boleh bernilai 'Ya', 'Belum', atau 'Sedang Tahap Penawaran'"
				}
			}

			for _, item := range getVisitList {
				formKey := item.Key
				// nameKey := item.Name
				valueBytes := c.Context().FormValue(formKey)
				valueStr := string(valueBytes)

				if formKey == "skul_id" {
					if getAccount.IsSkulid == nil || *getAccount.IsSkulid == 0 {
						if statusSkulid != "Ya" {
							continue
						}
					}
				}

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
						demoDescription := c.FormValue("demo_description")
						if demoDescription == "" {
							errors := map[string]string{
								"demo_description": "Deskripsi presentasi / demo harus diisi",
							}
							response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
							return c.Status(fiber.StatusBadRequest).JSON(response)
						}

						detailVisit[formKey+"_description"] = demoDescription

						demoDocumentation := c.FormValue("demo_documentation", "")
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
						detailVisit[formKey] = filePath
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
						bakFile := c.FormValue("bak_file", "")
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
						detailVisit[formKey] = filePath

						event_name := c.FormValue("event_name")
						if event_name == "" {
							errors := map[string]string{
								"event_name": "Nama event harus diisi",
							}
							response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
							return c.Status(fiber.StatusBadRequest).JSON(response)
						}

						detailVisit["event_name"] = helper.UppercaseTrim(event_name)

						amountDealingStr := c.FormValue("amount_dealing")
						if amountDealingStr == "" {
							errors := map[string]string{
								"amount_dealing": "Jumlah dealing harus diisi",
							}
							response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
							return c.Status(fiber.StatusBadRequest).JSON(response)
						}

						// Validasi bahwa amount_dealing adalah angka
						amountDealing, err := strconv.Atoi(amountDealingStr)
						if err != nil {
							errors := map[string]string{
								"amount_dealing": "Jumlah dealing harus berupa angka",
							}
							response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
							return c.Status(fiber.StatusBadRequest).JSON(response)
						}

						detailVisit["amount_dealing"] = strconv.Itoa(amountDealing)

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
				if (getAccount.IsSkulid != nil && *getAccount.IsSkulid == 1) || (statusSkulid != "" && statusSkulid == "Ya") {
					if formKey == "skul_id" {
						if valueStr == "1" {
							skulidDescription := c.FormValue("skulid_description")
							if skulidDescription == "" {
								errors := map[string]string{
									"skulid_description": "Deskripsi SkulID harus diisi",
								}
								response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
								return c.Status(fiber.StatusBadRequest).JSON(response)
							}

							detailVisit[formKey+"_description"] = skulidDescription
						}
					}
				}

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
		// check if other user clock in
		if req.Type == "Visit Account" {
			subjectIDStr := c.FormValue("subject_id")
			parsedSubjectID, _ := strconv.Atoi(subjectIDStr)
			subjectID = parsedSubjectID
			existingAbsenceUser, message, _ := h.absenceUserService.GetAbsenceUserToday(
				false,
				0,
				&req.Type,
				type_checking,
				actionType,
				"App\\Models\\Account",
				subjectID,
			)
			if existingAbsenceUser != nil {
				errors := map[string]string{
					"message": message,
				}
				response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
				return c.Status(fiber.StatusBadRequest).JSON(response)
			}
		}

		existingAbsenceUser, message, _ := h.absenceUserService.GetAbsenceUserToday(
			false,
			userID,
			&req.Type,
			type_checking,
			actionType,
			"",
			0,
		)

		typeClockIn := c.FormValue("type_clock_in")
		var evidenceImage *string
		status := 1

		if typeClockIn == "image" {
			status = 0
			notifyUser = true
			typeClockIn := c.FormValue("evidence_image", "")
			if typeClockIn == "" {
				errors := map[string]string{
					"evidence_image": "Dokumentasi demo harus diisi",
				}
				response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
				return c.Status(fiber.StatusBadRequest).JSON(response)
			}

			// Simpan gambar base64
			filePath, err := helper.SaveValidatedBase64File(typeClockIn, "public/uploads/evidence_absence")
			if err != nil {
				errors := map[string]string{
					"evidence_image": err.Error(),
				}
				response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
				return c.Status(fiber.StatusBadRequest).JSON(response)
			}

			// Simpan relative path
			evidenceImage = &filePath
		}

		if existingAbsenceUser != nil {
			errors := map[string]string{
				"message": message,
			}
			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}

		statusUint := uint(status)
		AbsenceUser, err := h.absenceUserService.CreateAbsenceUser(userID, subjectTypeStr, subjectID, &description, &req.Type, &req.Latitude, &req.Longitude, &statusUint, evidenceImage)
		if err != nil {
			response := helper.APIResponse(err.Error(), fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}

		if notifyUser {
			subjectIDStr := c.FormValue("subject_id")
			parsedSubjectID, _ := strconv.Atoi(subjectIDStr)
			getAccount, err := h.accountService.FindByAccountID(uint(parsedSubjectID), userRole, uint(territoryID), uint(userID))
			if err != nil {
				response := helper.APIResponse("Failed to fetch account", fiber.StatusInternalServerError, "error", nil)
				return c.Status(fiber.StatusInternalServerError).JSON(response)
			}

			getUser, err := h.userService.GetUserByID(uint(userID))
			if err != nil {
				response := helper.APIResponse("Failed to fetch user", fiber.StatusInternalServerError, "error", nil)
				return c.Status(fiber.StatusInternalServerError).JSON(response)
			}

			requestBody := map[string]string{
				"title":        "Approval Absensi Visit",
				"description":  fmt.Sprintf("Permintaan approval absensi visit account %s dari %s.", *getAccount.AccountName, getUser.Name),
				"callback_url": fmt.Sprintf("/visits?type=detail&id=%d", AbsenceUser.ID),
				"subject_type": "App\\Models\\AbsenceUser",
				"subject_id":   fmt.Sprintf("%d", AbsenceUser.ID),
			}
			_ = h.notificationService.CreateNotification(requestBody, []string{"Branch"}, userRole, territoryID, 0)

			requestBody = map[string]string{
				"message":      fmt.Sprintf("Permintaan approval absensi visit account %s dari %s.", *getAccount.AccountName, getUser.Name),
				"callback_url": fmt.Sprintf("/visits?type=detail&id=%d", AbsenceUser.ID),
			}
			err = h.smsSenderService.CreateSms(requestBody, []string{"Branch"}, userRole, territoryID, 0)

			fmt.Println("err", err)
		}

		response := helper.APIResponse(successMessage, successCode, "success", AbsenceUser)
		return c.Status(successCode).JSON(response)

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

		if req.Type == "Daily" {
			// Count different time
			duration := time.Since(existingAbsenceUser.ClockIn)

			// If less than 6 hour
			if duration < 6*time.Hour {
				response := helper.APIResponse("Absensi Harian belum mencapai 6 jam", fiber.StatusBadRequest, "error", nil)
				return c.Status(fiber.StatusBadRequest).JSON(response)
			}
		}

		status := 1
		statusUint := uint(status)
		AbsenceUser, err := h.absenceUserService.UpdateAbsenceUser(int(existingAbsenceUser.ID), userID, subjectTypeStr, subjectID, &description, &req.Type, &statusUint)
		if err != nil {
			response := helper.APIResponse(err.Error(), fiber.StatusUnauthorized, "error", AbsenceUser)
			return c.Status(fiber.StatusUnauthorized).JSON(response)
		}

		if req.Type == "Visit Account" {
			subjectIDStr := c.FormValue("subject_id")
			parsedSubjectID, _ := strconv.Atoi(subjectIDStr)
			getAccount, err := h.accountService.FindByAccountID(uint(parsedSubjectID), userRole, uint(territoryID), uint(userID))
			if err != nil {
				response := helper.APIResponse("Failed to fetch account", fiber.StatusInternalServerError, "error", nil)
				return c.Status(fiber.StatusInternalServerError).JSON(response)
			}

			if getAccount.IsSkulid == nil || *getAccount.IsSkulid != 1 {
				isSkulid := c.FormValue("status_skulid")
				statusSkulid := 0

				if isSkulid == "Ya" {
					statusSkulid = 1
				} else if isSkulid == "Belum" {
					statusSkulid = 0
				} else if isSkulid == "Sedang Tahap Penawaran" {
					statusSkulid = 2
				}
				requestBody := map[string]interface{}{
					"is_skulid": statusSkulid,
				}
				err := h.accountService.UpdateFields(uint(parsedSubjectID), requestBody)
				if err != nil {
					response := helper.APIResponse(err.Error(), fiber.StatusInternalServerError, "error", nil)
					return c.Status(fiber.StatusInternalServerError).JSON(response)
				}
			}

			VisitHistory, err := h.visitHistoryService.CreateVisitHistory(userID, subjectTypeStr, subjectID, int(AbsenceUser.ID), kpiYae, &description, detailVisit)
			if err != nil {
				response := helper.APIResponse(err.Error(), fiber.StatusUnauthorized, "error", VisitHistory)
				return c.Status(fiber.StatusUnauthorized).JSON(response)
			}
		}

		response := helper.APIResponse(successMessage, successCode, "success", AbsenceUser)
		return c.Status(successCode).JSON(response)
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

func (h *AbsenceUserHandler) ExportRawAbsenceUsers(c *fiber.Ctx) error {
	// Ambil user_id dari query parameter atau dari locals
	var userID int
	if queryUserID := c.Query("user_id"); queryUserID != "" {
		parsedID, err := strconv.Atoi(queryUserID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "user_id tidak valid",
			})
		}
		userID = parsedID
	} else {
		local := c.Locals("user_id")
		if local == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "user_id tidak ditemukan",
			})
		}
		userID = local.(int)
	}

	// Ambil semua filter dari query string
	filters := map[string]string{
		"search":     c.Query("search", ""),
		"order_by":   c.Query("order_by", "id"),
		"order":      c.Query("order", "DESC"),
		"start_date": c.Query("start_date", ""),
		"end_date":   c.Query("end_date", ""),
		"status":     c.Query("status", "1"),
		"all_user":   c.Query("all_user", "0"),
	}

	// Parse tambahan parameter jika perlu
	month, _ := strconv.Atoi(c.Query("month", "0"))
	year, _ := strconv.Atoi(c.Query("year", "0"))
	absenceType := c.Query("type", "")

	// Panggil service untuk generate file Excel
	userRole := c.Locals("user_role").(string)
	territoryID := c.Locals("territory_id").(int)
	filtersUser := map[string]string{
		"search":      "",
		"order_by":    "id",
		"order":       "DESC",
		"start_date":  "",
		"end_date":    "",
		"user_status": "active",
	}
	getAllUsers, _, err := h.userService.GetAllUsers(0, false, 1, filtersUser, []string{"YAE"}, false, userRole, territoryID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data user",
		})
	}

	userIDs := []int{}
	for _, user := range getAllUsers {
		userIDs = append(userIDs, int(user.ID)) // konversi uint ke int
	}

	excelFile, err := h.absenceUserService.GenerateAbsenceExcel(userID, filters, month, year, absenceType, userRole, territoryID, userIDs)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal membuat file Excel",
		})
	}

	// Set header sebagai file response
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename=absence_export.xlsx")

	// Kirim stream file
	return c.SendStream(excelFile)
}

func (h *AbsenceUserHandler) ExportResumeMonthlyAbsenceUsers(c *fiber.Ctx) error {
	// Ambil user_id dari query parameter atau dari locals
	var userID int
	if queryUserID := c.Query("user_id"); queryUserID != "" {
		parsedID, err := strconv.Atoi(queryUserID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "user_id tidak valid",
			})
		}
		userID = parsedID
	} else {
		local := c.Locals("user_id")
		if local == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "user_id tidak ditemukan",
			})
		}
		userID = local.(int)
	}

	// Ambil semua filter dari query string
	filters := map[string]string{
		"search":     c.Query("search", ""),
		"order_by":   c.Query("order_by", "id"),
		"order":      c.Query("order", "DESC"),
		"start_date": c.Query("start_date", ""),
		"end_date":   c.Query("end_date", ""),
		"status":     c.Query("status", "1"),
		"all_user":   c.Query("all_user", "0"),
	}

	// Parse tambahan parameter jika perlu
	month, _ := strconv.Atoi(c.Query("month", "0"))
	year, _ := strconv.Atoi(c.Query("year", "0"))
	absenceType := c.Query("type", "")
	userRole := c.Locals("user_role").(string)
	territoryID := c.Locals("territory_id").(int)

	// Panggil service untuk generate file Excel
	excelFile, err := h.absenceUserService.GenerateAbsenceResumeExcel(userID, filters, month, year, absenceType, userRole, territoryID, []int{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal membuat file Excel",
		})
	}

	// Set header sebagai file response
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename=absence_resume_export.xlsx")

	// Kirim stream file
	return c.SendStream(excelFile)
}

func (h *AbsenceUserHandler) HandleAbsenceApproval(c *fiber.Ctx) error {
	userRole := c.Locals("user_role").(string)
	territoryID := c.Locals("territory_id").(int)
	userID := c.Locals("user_id").(int)
	id := c.Params("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		errors := map[string]string{
			"id": "ID tidak valid",
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	status := c.Params("status")
	if status != "accept" && status != "reject" {
		errors := map[string]string{
			"status": "Status tidak valid",
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	AbsenceUser, err := h.absenceUserService.GetAbsenceUserByID(intID)
	if err != nil {
		response := helper.APIResponse("Failed to fetch Absence User", fiber.StatusNotFound, "error", nil)
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	if *AbsenceUser.Status != 0 {
		response := helper.APIResponse("Absence already accepted", fiber.StatusOK, "error", nil)
		return c.Status(fiber.StatusOK).JSON(response)
	}

	if err != nil {
		response := helper.APIResponse("Invalid status value", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if status == "accept" {
		var err error
		updateData := map[string]interface{}{
			"status": 1,
		}

		err = h.absenceUserService.UpdateFields(uint(intID), updateData)
		if err != nil {
			response := helper.APIResponse(err.Error(), fiber.StatusUnauthorized, "error", AbsenceUser)
			return c.Status(fiber.StatusUnauthorized).JSON(response)
		}
		var statusInt uint = 1
		AbsenceUser.Status = &statusInt

		getAccount, err := h.accountService.FindByAccountID(uint(*AbsenceUser.SubjectID), userRole, uint(territoryID), uint(userID))
		if err != nil {
			response := helper.APIResponse("Failed to fetch account", fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}

		requestBody := map[string]string{
			"title":        "Pengajuan Visit Diterima",
			"description":  fmt.Sprintf("Permintaan approval absensi visit account %s anda diterima.", *getAccount.AccountName),
			"callback_url": fmt.Sprintf("/my-task/%d", *AbsenceUser.SubjectID),
			"subject_type": "App\\Models\\AbsenceUser",
			"subject_id":   fmt.Sprintf("%d", AbsenceUser.ID),
		}
		_ = h.notificationService.CreateNotification(requestBody, []string{}, userRole, territoryID, int(*AbsenceUser.UserID))

		requestBody = map[string]string{
			"message":      fmt.Sprintf("Permintaan approval absensi visit account %s anda diterima.", *getAccount.AccountName),
			"callback_url": fmt.Sprintf("/my-task/%d", *AbsenceUser.SubjectID),
		}
		_ = h.smsSenderService.CreateSms(requestBody, []string{}, userRole, territoryID, int(*AbsenceUser.UserID))

	} else if status == "reject" {
		err = h.absenceUserService.DeleteAbsenceUser(intID)
		if err != nil {
			response := helper.APIResponse("Failed to reject", fiber.StatusOK, "error", nil)
			return c.Status(fiber.StatusOK).JSON(response)
		}

		if AbsenceUser.EvidenceImage != nil && *AbsenceUser.EvidenceImage != "" {
			err := helper.DeleteFileIfExists(*AbsenceUser.EvidenceImage)
			if err != nil {
				response := helper.APIResponse("Failed to reject", fiber.StatusOK, "error", err.Error())
				return c.Status(fiber.StatusOK).JSON(response)
			}
		}

		getAccount, err := h.accountService.FindByAccountID(uint(*AbsenceUser.SubjectID), userRole, uint(territoryID), uint(userID))
		if err != nil {
			response := helper.APIResponse("Failed to fetch account", fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}

		requestBody := map[string]string{
			"title":        "Pengajuan Visit Ditolak",
			"description":  fmt.Sprintf("Permintaan approval absensi visit account %s anda ditolak. Silahkan ajukan ulang.", *getAccount.AccountName),
			"callback_url": fmt.Sprintf("/my-task/%d", *AbsenceUser.SubjectID),
			"subject_type": "App\\Models\\AbsenceUser",
			"subject_id":   fmt.Sprintf("%d", AbsenceUser.ID),
		}
		_ = h.notificationService.CreateNotification(requestBody, []string{}, userRole, territoryID, userID)

		requestBody = map[string]string{
			"message":      fmt.Sprintf("Permintaan approval absensi visit account %s anda ditolak. Silahkan ajukan ulang.", *getAccount.AccountName),
			"callback_url": fmt.Sprintf("/my-task/%d", *AbsenceUser.SubjectID),
		}
		_ = h.smsSenderService.CreateSms(requestBody, []string{}, userRole, territoryID, userID)
	}

	// Return response
	responseData := map[string]interface{}{
		"absence": AbsenceUser,
	}

	response := helper.APIResponse("Approval Absence User Successfully", fiber.StatusOK, "success", responseData)
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
