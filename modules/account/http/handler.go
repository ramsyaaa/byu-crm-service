package http

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"byu-crm-service/helper"
	"byu-crm-service/models"
	"byu-crm-service/modules/account/service"
	"byu-crm-service/modules/account/validation"

	absenceUserService "byu-crm-service/modules/absence-user/service"
	accountFacultyService "byu-crm-service/modules/account-faculty/service"
	accountMemberService "byu-crm-service/modules/account-member/service"
	accountScheduleService "byu-crm-service/modules/account-schedule/service"
	accountTypeCampusDetailService "byu-crm-service/modules/account-type-campus-detail/service"
	accountTypeCommunityDetailService "byu-crm-service/modules/account-type-community-detail/service"
	accountTypeSchoolDetailService "byu-crm-service/modules/account-type-school-detail/service"
	contactAccountService "byu-crm-service/modules/contact-account/service"
	productService "byu-crm-service/modules/product/service"
	socialMediaService "byu-crm-service/modules/social-media/service"

	"github.com/gofiber/fiber/v2"
)

type AccountHandler struct {
	service                           service.AccountService
	contactAccountService             contactAccountService.ContactAccountService
	socialMediaService                socialMediaService.SocialMediaService
	accountTypeSchoolDetailService    accountTypeSchoolDetailService.AccountTypeSchoolDetailService
	accountFacultyService             accountFacultyService.AccountFacultyService
	accountMemberService              accountMemberService.AccountMemberService
	accountScheduleService            accountScheduleService.AccountScheduleService
	accountTypeCampusDetailService    accountTypeCampusDetailService.AccountTypeCampusDetailService
	accountTypeCommunityDetailService accountTypeCommunityDetailService.AccountTypeCommunityDetailService
	productService                    productService.ProductService
	absenceUserService                absenceUserService.AbsenceUserService
}

func NewAccountHandler(
	service service.AccountService,
	contactAccountService contactAccountService.ContactAccountService,
	socialMediaService socialMediaService.SocialMediaService,
	accountTypeSchoolDetailService accountTypeSchoolDetailService.AccountTypeSchoolDetailService,
	accountFacultyService accountFacultyService.AccountFacultyService,
	accountMemberService accountMemberService.AccountMemberService,
	accountScheduleService accountScheduleService.AccountScheduleService,
	accountTypeCampusDetailService accountTypeCampusDetailService.AccountTypeCampusDetailService,
	accountTypeCommunityDetailService accountTypeCommunityDetailService.AccountTypeCommunityDetailService,
	productService productService.ProductService,
	absenceUserService absenceUserService.AbsenceUserService) *AccountHandler {

	return &AccountHandler{
		service:                           service,
		contactAccountService:             contactAccountService,
		socialMediaService:                socialMediaService,
		accountTypeSchoolDetailService:    accountTypeSchoolDetailService,
		accountFacultyService:             accountFacultyService,
		accountMemberService:              accountMemberService,
		accountScheduleService:            accountScheduleService,
		accountTypeCampusDetailService:    accountTypeCampusDetailService,
		accountTypeCommunityDetailService: accountTypeCommunityDetailService,
		productService:                    productService,
		absenceUserService:                absenceUserService}
}

func (h *AccountHandler) GetAllAccounts(c *fiber.Ctx) error {
	// Default query params
	filters := map[string]string{
		"search":           c.Query("search", ""),
		"order_by":         c.Query("order_by", "id"),
		"order":            c.Query("order", "DESC"),
		"start_date":       c.Query("start_date", ""),
		"end_date":         c.Query("end_date", ""),
		"account_category": c.Query("account_category", ""),
		"account_type":     c.Query("account_type", ""),
		"only_skulid":      c.Query("only_skulid", "0"),
		"is_priority":      c.Query("is_priority", "0"),
		"priority":         c.Query("priority", "P1"),
	}

	// Parse integer and boolean values
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	paginate, _ := strconv.ParseBool(c.Query("paginate", "true"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	userRole := c.Locals("user_role").(string)
	territoryID := c.Locals("territory_id").(int)
	userID := c.Locals("user_id").(int)
	onlyUserPic, _ := strconv.ParseBool(c.Query("only_user_pic", "0"))
	excludeVisited, _ := strconv.ParseBool(c.Query("exclude_visited", "false"))

	// Call service with filters
	accounts, total, err := h.service.GetAllAccounts(limit, paginate, page, filters, userRole, territoryID, userID, onlyUserPic, excludeVisited)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch accounts",
			"error":   err.Error(),
		})
	}

	// Return response
	responseData := map[string]interface{}{
		"accounts": accounts,
		"total":    total,
		"page":     page,
	}

	response := helper.APIResponse("Get Accounts Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *AccountHandler) GetCountAccount(c *fiber.Ctx) error {
	var territoryID int
	var userRole string
	var err error

	// Ambil user_role dari query jika ada, jika tidak ambil dari locals
	userRoleParam := c.Query("user_role")
	if userRoleParam != "" {
		userRole = userRoleParam
	} else {
		var ok bool
		userRole, ok = c.Locals("user_role").(string)
		if !ok {
			response := helper.APIResponse("Unauthorized: Invalid user role", fiber.StatusUnauthorized, "error", nil)
			return c.Status(fiber.StatusUnauthorized).JSON(response)
		}
	}

	// Ambil territory_id dari query jika ada, jika tidak ambil dari locals
	territoryIDParam := c.Query("territory_id")
	if territoryIDParam != "" {
		territoryID, err = strconv.Atoi(territoryIDParam)
		if err != nil {
			response := helper.APIResponse("Bad Request: Invalid territory_id", fiber.StatusBadRequest, "error", nil)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	} else {
		var ok bool
		territoryID, ok = c.Locals("territory_id").(int)
		if !ok {
			response := helper.APIResponse("Unauthorized: Invalid territory ID", fiber.StatusUnauthorized, "error", nil)
			return c.Status(fiber.StatusUnauthorized).JSON(response)
		}
	}

	// Call service with filters
	total, categories, territories, territory_info, err := h.service.CountAccount(userRole, territoryID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch accounts",
			"error":   err.Error(),
		})
	}

	// Return response
	responseData := map[string]interface{}{
		"total":          total,
		"territory_info": territory_info,
		"categories":     categories,
		"territories":    territories,
	}

	response := helper.APIResponse("Get Accounts Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *AccountHandler) CheckAlreadyUpdateData(c *fiber.Ctx) error {
	idParam := c.Params("id")
	userID := c.Locals("user_id").(int)
	// Convert to int
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response := helper.APIResponse("Invalid ID format", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Ambil territory_id dari query jika ada, jika tidak ambil dari locals
	clockInStr := c.Query("clock_in")
	if clockInStr == "" {
		response := helper.APIResponse("Parameter 'clock_in' harus diisi", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	clockInTime, err := time.Parse("2006-01-02", clockInStr)
	if err != nil {
		response := helper.APIResponse("Format tanggal 'clock_in' tidak valid. Gunakan format YYYY-MM-DD", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Call service with filters
	status, err := h.service.CheckAlreadyUpdateData(id, clockInTime, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch data update",
			"error":   err.Error(),
		})
	}

	// Return response
	responseData := map[string]interface{}{
		"status": status,
	}

	response := helper.APIResponse("Check Data Update Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *AccountHandler) GetAccountVisitCounts(c *fiber.Ctx) error {
	// Default query params
	filters := map[string]string{
		"search":     c.Query("search", ""),
		"order_by":   c.Query("order_by", "id"),
		"order":      c.Query("order", "DESC"),
		"start_date": c.Query("start_date", ""),
		"end_date":   c.Query("end_date", ""),
	}

	// Parse integer and boolean values
	userRole := c.Locals("user_role").(string)
	territoryID := c.Locals("territory_id").(int)
	userID := c.Locals("user_id").(int)

	// Call service with filters
	visited_account, not_visited, total, err := h.service.GetAccountVisitCounts(filters, userRole, territoryID, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch accounts",
			"error":   err.Error(),
		})
	}

	// Return response
	responseData := map[string]interface{}{
		"visited_account": visited_account,
		"not_visited":     not_visited,
		"total":           total,
	}

	response := helper.APIResponse("Get Accounts Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *AccountHandler) GetAccountById(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	// Get id from param
	idParam := c.Params("id")
	userRole := c.Locals("user_role").(string)
	territoryID := c.Locals("territory_id").(int)

	// Convert to int
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response := helper.APIResponse("Invalid ID format", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	account, err := h.service.FindByAccountID(uint(id), userRole, uint(territoryID), uint(userID))
	if err != nil {
		response := helper.APIResponse("Account not found", fiber.StatusNotFound, "error", nil)
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	response := helper.APIResponse("Success get account", fiber.StatusOK, "success", account)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *AccountHandler) CreateAccount(c *fiber.Ctx) error {
	// Add a timeout context to prevent long-running operations
	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	// Use a recovery function to catch any panics
	defer func() {
		if r := recover(); r != nil {
			helper.LogError(c, fmt.Sprintf("Panic in Create Account: %v", r))
			response := helper.APIResponse("Internal server error", fiber.StatusInternalServerError, "error", r)
			c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}()

	// Get user information from context
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		response := helper.APIResponse("Unauthorized: Invalid user ID", fiber.StatusUnauthorized, "error", nil)
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	territoryID, ok := c.Locals("territory_id").(int)
	if !ok {
		response := helper.APIResponse("Unauthorized: Invalid territory ID", fiber.StatusUnauthorized, "error", nil)
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	userRole, ok := c.Locals("user_role").(string)
	if !ok {
		response := helper.APIResponse("Unauthorized: Invalid user role", fiber.StatusUnauthorized, "error", nil)
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	// Parse request body with error handling
	req := new(validation.ValidateRequest)

	if err := c.BodyParser(req); err != nil {
		// Check for specific EOF error
		if err.Error() == "unexpected EOF" {
			response := helper.APIResponse("Invalid request: Unexpected end of JSON input", fiber.StatusBadRequest, "error", nil)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}

		response := helper.APIResponse("Invalid request format: "+err.Error(), fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	// Request Validation with context
	select {
	case <-ctx.Done():
		response := helper.APIResponse("Request timeout during validation", fiber.StatusRequestTimeout, "error", nil)
		return c.Status(fiber.StatusRequestTimeout).JSON(response)
	default:
		errors := validation.ValidateCreate(req)
		if errors != nil {
			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	}

	if req.AccountCategory != "" && req.AccountCategory == "SEKOLAH" {

		errors := validation.ValidateSchool(req, true, 0, userRole, territoryID, userID)
		if errors != nil {
			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	} else if req.AccountCategory != "" && req.AccountCategory == "KAMPUS" {

		errors := validation.ValidateCampus(req, true, 0, userRole, territoryID, userID)
		if errors != nil {
			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	} else if req.AccountCategory != "" && req.AccountCategory == "KOMUNITAS" {

		errors := validation.ValidateCommunity(req, true, 0, userRole, territoryID, userID)
		if errors != nil {
			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	}

	// Create Account with context and error handling
	reqMap := make(map[string]interface{})

	// Marshal request to JSON with timeout
	var reqBytes []byte
	var marshalErr error

	select {
	case <-ctx.Done():
		response := helper.APIResponse("Request timeout during marshaling", fiber.StatusRequestTimeout, "error", nil)
		return c.Status(fiber.StatusRequestTimeout).JSON(response)
	default:
		reqBytes, marshalErr = json.Marshal(req)
		if marshalErr != nil {
			helper.LogError(c, fmt.Sprintf("Failed to marshal request: %v", marshalErr))
			response := helper.APIResponse("Failed to process request data", fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	// Unmarshal JSON to map with timeout
	var unmarshalErr error

	select {
	case <-ctx.Done():
		response := helper.APIResponse("Request timeout during unmarshaling", fiber.StatusRequestTimeout, "error", nil)
		return c.Status(fiber.StatusRequestTimeout).JSON(response)
	default:
		unmarshalErr = json.Unmarshal(reqBytes, &reqMap)
		if unmarshalErr != nil {
			helper.LogError(c, fmt.Sprintf("Failed to unmarshal request: %v", unmarshalErr))
			response := helper.APIResponse("Failed to process request data", fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	// Call service with timeout
	var account []models.Account
	var serviceErr error

	select {
	case <-ctx.Done():
		response := helper.APIResponse("Request timeout during account creation", fiber.StatusRequestTimeout, "error", nil)
		return c.Status(fiber.StatusRequestTimeout).JSON(response)
	default:
		account, serviceErr = h.service.CreateAccount(reqMap, userID)
		if serviceErr != nil {
			helper.LogError(c, fmt.Sprintf("Failed to create account: %v", serviceErr))
			response := helper.APIResponse("Failed to create account: "+serviceErr.Error(), fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	_, _ = h.contactAccountService.InsertContactAccount(reqMap, account[0].ID)
	_, _ = h.socialMediaService.InsertSocialMedia(reqMap, "App\\Models\\Account", account[0].ID)
	_, _ = h.productService.InsertProductAccount(reqMap, account[0].ID)

	// Safely handle account_category - check if it exists and is not nil
	if category, exists := reqMap["account_category"]; exists && category != nil {
		// Try to convert to string, with proper type assertion check
		categoryStr, ok := category.(string)
		if !ok {
			// Log the error but don't panic
			helper.LogError(c, fmt.Sprintf("account_category is not a string: %v (type: %T)", category, category))
		} else {
			switch categoryStr {
			case "SEKOLAH":
				_, _ = h.accountTypeSchoolDetailService.Insert(reqMap, account[0].ID)
			case "KAMPUS":
				_, _ = h.accountFacultyService.Insert(reqMap, account[0].ID)
				_, _ = h.accountMemberService.Insert(reqMap, "App\\Models\\Account", account[0].ID, "year", "amount")
				_, _ = h.accountMemberService.Insert(reqMap, "App\\Models\\AccountLecture", account[0].ID, "year_lecture", "amount_lecture")
				_, _ = h.accountScheduleService.Insert(reqMap, "App\\Models\\Account", account[0].ID)
				_, _ = h.accountTypeCampusDetailService.Insert(reqMap, account[0].ID)
			case "KOMUNITAS":
				_, _ = h.accountMemberService.Insert(reqMap, "App\\Models\\Account", account[0].ID, "year", "amount")
				_, _ = h.accountTypeCommunityDetailService.Insert(reqMap, account[0].ID)
				_, _ = h.accountScheduleService.Insert(reqMap, "App\\Models\\Account", account[0].ID)
			default:
				// Log unexpected category value
				helper.LogError(c, fmt.Sprintf("Unexpected account_category value: %s", categoryStr))
			}
		}
	} else if !exists {
		// Log missing account_category
		helper.LogError(c, "account_category field is missing in the request")
	} else {
		// Log nil account_category
		helper.LogError(c, "account_category field is nil")
	}

	// Return success response
	response := helper.APIResponse("Create Account Succsesfully", fiber.StatusOK, "success", account[0])
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *AccountHandler) UpdateAccount(c *fiber.Ctx) error {
	// Add a timeout context to prevent long-running operations
	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	// Use a recovery function to catch any panics
	defer func() {
		if r := recover(); r != nil {
			helper.LogError(c, fmt.Sprintf("Panic in Update Account: %v", r))
			response := helper.APIResponse("Internal server error", fiber.StatusInternalServerError, "error", r)
			c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}()

	// Get user information from context
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		response := helper.APIResponse("Unauthorized: Invalid user ID", fiber.StatusUnauthorized, "error", nil)
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	territoryID, ok := c.Locals("territory_id").(int)
	if !ok {
		response := helper.APIResponse("Unauthorized: Invalid territory ID", fiber.StatusUnauthorized, "error", nil)
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	userRole, ok := c.Locals("user_role").(string)
	if !ok {
		response := helper.APIResponse("Unauthorized: Invalid user role", fiber.StatusUnauthorized, "error", nil)
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	// Get and validate account ID
	accountIDStr := c.Params("id")
	if accountIDStr == "" {
		response := helper.APIResponse("Account ID is required", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		response := helper.APIResponse("Invalid Account ID", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Parse request body with error handling
	req := new(validation.ValidateRequest)
	if err := c.BodyParser(req); err != nil {
		// Check for specific EOF error
		if err.Error() == "unexpected EOF" {
			response := helper.APIResponse("Invalid request: Unexpected end of JSON input", fiber.StatusBadRequest, "error", nil)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}

		response := helper.APIResponse("Invalid request format: "+err.Error(), fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Request Validation with context
	select {
	case <-ctx.Done():
		response := helper.APIResponse("Request timeout during validation", fiber.StatusRequestTimeout, "error", nil)
		return c.Status(fiber.StatusRequestTimeout).JSON(response)
	default:
		errors := validation.ValidateCreate(req)
		if errors != nil {
			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	}

	if req.AccountCategory != "" && req.AccountCategory == "SEKOLAH" {
		errors := validation.ValidateSchool(req, false, accountID, userRole, territoryID, userID)
		if errors != nil {
			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	} else if req.AccountCategory != "" && req.AccountCategory == "KAMPUS" {

		errors := validation.ValidateCampus(req, false, accountID, userRole, territoryID, userID)
		if errors != nil {
			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	} else if req.AccountCategory != "" && req.AccountCategory == "KOMUNITAS" {

		errors := validation.ValidateCommunity(req, false, accountID, userRole, territoryID, userID)
		if errors != nil {
			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	}

	// Update Account with context and error handling
	reqMap := make(map[string]interface{})

	// Marshal request to JSON with timeout
	var reqBytes []byte
	var marshalErr error

	select {
	case <-ctx.Done():
		response := helper.APIResponse("Request timeout during marshaling", fiber.StatusRequestTimeout, "error", nil)
		return c.Status(fiber.StatusRequestTimeout).JSON(response)
	default:
		reqBytes, marshalErr = json.Marshal(req)
		if marshalErr != nil {
			helper.LogError(c, fmt.Sprintf("Failed to marshal request: %v", marshalErr))
			response := helper.APIResponse("Failed to process request data", fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	// Unmarshal JSON to map with timeout
	var unmarshalErr error

	select {
	case <-ctx.Done():
		response := helper.APIResponse("Request timeout during unmarshaling", fiber.StatusRequestTimeout, "error", nil)
		return c.Status(fiber.StatusRequestTimeout).JSON(response)
	default:
		unmarshalErr = json.Unmarshal(reqBytes, &reqMap)
		if unmarshalErr != nil {
			helper.LogError(c, fmt.Sprintf("Failed to unmarshal request: %v", unmarshalErr))
			response := helper.APIResponse("Failed to process request data", fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	// Call service with timeout
	var account []models.Account
	var serviceErr error

	select {
	case <-ctx.Done():
		response := helper.APIResponse("Request timeout during account update", fiber.StatusRequestTimeout, "error", nil)
		return c.Status(fiber.StatusRequestTimeout).JSON(response)
	default:
		account, serviceErr = h.service.UpdateAccount(reqMap, accountID, userRole, territoryID, userID)
		if serviceErr != nil {
			helper.LogError(c, fmt.Sprintf("Failed to update account: %v", serviceErr))
			response := helper.APIResponse("Failed to update account: "+serviceErr.Error(), fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	visitAccount := "Visit Account"
	existingAbsenceUser, _, _ := h.absenceUserService.GetAbsenceUserToday(
		false,
		userID,
		&visitAccount,
		"monthly",
		"Clock In",
		"App\\Models\\Account",
		accountID,
	)

	if existingAbsenceUser != nil {
		var subjectType *string = nil
		var subjectID *uint = nil

		err := h.service.CreateHistoryActivityAccount(uint(userID), uint(accountID), "Update Account", subjectType, subjectID)
		if err != nil {
			log.Println("Failed to create update history:", err)
		}
	}

	_, _ = h.contactAccountService.InsertContactAccount(reqMap, account[0].ID)
	_, _ = h.socialMediaService.InsertSocialMedia(reqMap, "App\\Models\\Account", account[0].ID)
	_, _ = h.productService.InsertProductAccount(reqMap, account[0].ID)

	// Safely handle account_category - check if it exists and is not nil
	if category, exists := reqMap["account_category"]; exists && category != nil {
		// Try to convert to string, with proper type assertion check
		categoryStr, ok := category.(string)
		if !ok {
			// Log the error but don't panic
			helper.LogError(c, fmt.Sprintf("account_category is not a string: %v (type: %T)", category, category))
		} else {
			switch categoryStr {
			case "SEKOLAH":
				_, _ = h.accountTypeSchoolDetailService.Insert(reqMap, account[0].ID)
			case "KAMPUS":
				_, _ = h.accountFacultyService.Insert(reqMap, account[0].ID)
				_, _ = h.accountMemberService.Insert(reqMap, "App\\Models\\Account", account[0].ID, "year", "amount")
				_, _ = h.accountMemberService.Insert(reqMap, "App\\Models\\AccountLecture", account[0].ID, "year_lecture", "amount_lecture")
				_, _ = h.accountScheduleService.Insert(reqMap, "App\\Models\\Account", account[0].ID)
				_, _ = h.accountTypeCampusDetailService.Insert(reqMap, account[0].ID)
			case "KOMUNITAS":
				_, _ = h.accountTypeCommunityDetailService.Insert(reqMap, account[0].ID)
				_, _ = h.accountScheduleService.Insert(reqMap, "App\\Models\\Account", account[0].ID)
			default:
				// Log unexpected category value
				helper.LogError(c, fmt.Sprintf("Unexpected account_category value: %s", categoryStr))
			}
		}
	} else if !exists {
		// Log missing account_category
		helper.LogError(c, "account_category field is missing in the request")
	} else {
		// Log nil account_category
		helper.LogError(c, "account_category field is nil")
	}

	// Return success response
	response := helper.APIResponse("Update Account Succsesfully", fiber.StatusOK, "success", account[0])
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *AccountHandler) Import(c *fiber.Ctx) error {
	// Validate the uploaded file
	if err := validation.ValidateAccountRequest(c); err != nil {
		return err // Error response is already handled in the validation function
	}

	// Get the file from the request
	file, err := c.FormFile("file_csv")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File is required"})
	}

	tempPath := "./temp/" + file.Filename

	// Ensure the temp directory exists
	if _, err := os.Stat("./temp/"); os.IsNotExist(err) {
		if err := os.Mkdir("./temp/", os.ModePerm); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create temp directory"})
		}
	}

	if err := c.SaveFile(file, tempPath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save file"})
	}

	// Retrieve user_id from the form
	userID := c.FormValue("user_id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID is required"})
	}

	go func() {
		defer os.Remove(tempPath)

		// Process file with timeout
		_, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		f, err := os.Open(tempPath)
		if err != nil {
			fmt.Println("Failed to open file:", err)
			return
		}
		defer f.Close()

		reader := csv.NewReader(f)
		rows, _ := reader.ReadAll()

		for i, row := range rows {
			if i == 0 {
				continue // Skip header
			}
			if err := h.service.ProcessAccount(row); err != nil {
				fmt.Println("Error processing row:", err)
				return
			}
		}

		// Send notification
		notificationURL := os.Getenv("NOTIFICATION_URL") + "/api/notification/create"
		payload := map[string]interface{}{
			"model":    "App\\Models\\Account",
			"model_id": 0, // Replace with actual model ID if needed
			"user_id":  userID,
			"data": map[string]string{
				"title":        "Import Account",
				"description":  "Import Account",
				"callback_url": "/accounts",
			},
		}

		payloadBytes, _ := json.Marshal(payload)
		resp, err := http.Post(notificationURL, "application/json", bytes.NewReader(payloadBytes))
		if err != nil {
			fmt.Println("Failed to send notification:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println("Notification API responded with status:", resp.StatusCode)
		} else {
			var responseMap map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
				fmt.Println("Failed to decode response:", err)
				return
			}
			fmt.Println("Notification sent successfully:", responseMap["message"])
		}
	}()

	return c.JSON(fiber.Map{
		"message": "File upload successful, processing in background",
		"status":  "success",
	})
}

func (h *AccountHandler) UpdatePic(c *fiber.Ctx) error {
	territoryID := c.Locals("territory_id").(int)
	userRole := c.Locals("user_role").(string)
	userID := c.Locals("user_id").(int)
	// Ambil Account ID dari URL
	accountIDStr := c.Params("id")
	if accountIDStr == "" {
		response := helper.APIResponse("Account ID is required", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Convert Account ID to integer
	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		response := helper.APIResponse("Invalid Account ID", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Panggil service untuk update pic
	account, err := h.service.UpdatePic(accountID, userRole, territoryID, userID)
	if err != nil {
		response := helper.APIResponse("Failed to update pic", fiber.StatusInternalServerError, "error", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Return success response
	response := helper.APIResponse("PIC Account successfully updated", fiber.StatusOK, "success", account)
	return c.Status(fiber.StatusOK).JSON(response)
}

func saveFileToLocal(file *multipart.FileHeader, directory string, allowedFormats []string) (*string, error) {
	// Validate file type
	ext := filepath.Ext(file.Filename)
	ext = strings.ToLower(ext)

	// Check if the file extension is allowed
	isValidExt := false
	for _, allowedExt := range allowedFormats {
		if ext == allowedExt {
			isValidExt = true
			break
		}
	}

	if !isValidExt {
		return nil, fmt.Errorf("invalid file format. Allowed formats are: %v", allowedFormats)
	}

	// Define a unique file name
	filename := fmt.Sprintf("%s%s", generateUniqueID(), ext)

	// Define the full path where to save the file
	savePath := filepath.Join("public", directory, filename)

	// Create the directory if it doesn't exist
	dir := filepath.Dir(savePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, fmt.Errorf("failed to create directories for file storage: %v", err)
		}
	}

	// Open the file from the incoming multipart request
	fileSrc, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer fileSrc.Close()

	// Create the destination file on the server
	fileDest, err := os.Create(savePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %v", err)
	}
	defer fileDest.Close()

	// Copy the content of the uploaded file to the destination file
	_, err = io.Copy(fileDest, fileSrc)
	if err != nil {
		return nil, fmt.Errorf("failed to save file: %v", err)
	}

	// Return the relative path to the saved file
	filePath := fmt.Sprintf("/%s/%s", directory, filename)
	return &filePath, nil
}

func generateUniqueID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
