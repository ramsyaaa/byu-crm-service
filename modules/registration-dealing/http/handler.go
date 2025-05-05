package http

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"byu-crm-service/helper"
	"byu-crm-service/modules/registration-dealing/response"
	"byu-crm-service/modules/registration-dealing/service"
	"byu-crm-service/modules/registration-dealing/validation"

	"github.com/gofiber/fiber/v2"
)

type RegistrationDealingHandler struct {
	service service.RegistrationDealingService
}

func NewRegistrationDealingHandler(
	service service.RegistrationDealingService) *RegistrationDealingHandler {

	return &RegistrationDealingHandler{
		service: service}
}

func (h *RegistrationDealingHandler) GetAllRegistrationDealings(c *fiber.Ctx) error {
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
	accountID, _ := strconv.Atoi(c.Query("account_id", "0"))
	event_name := c.Query("event_name", "")

	// Call service with filters
	registrationDealings, total, err := h.service.GetAllRegistrationDealings(limit, paginate, page, filters, accountID, event_name)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to fetch registration dealing",
			"error":   err.Error(),
		})
	}

	// Return response
	responseData := map[string]interface{}{
		"registration_dealings": registrationDealings,
		"total":                 total,
		"page":                  page,
	}

	response := helper.APIResponse("Get Registration Dealing Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *RegistrationDealingHandler) GetRegistrationDealingById(c *fiber.Ctx) error {
	// Get id from param
	idParam := c.Params("id")

	// Convert to int
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response := helper.APIResponse("Invalid ID format", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	registrationDealing, err := h.service.FindByRegistrationDealingID(uint(id))
	if err != nil {
		response := helper.APIResponse("Registration dealing not found", fiber.StatusNotFound, "error", nil)
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	response := helper.APIResponse("Success get registration dealing", fiber.StatusOK, "success", registrationDealing)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *RegistrationDealingHandler) CreateRegistrationDealing(c *fiber.Ctx) error {
	// Add a timeout context to prevent long-running operations
	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	// Use a recovery function to catch any panics
	defer func() {
		if r := recover(); r != nil {
			helper.LogError(c, fmt.Sprintf("Panic in Create Registration Dealing: %v", r))
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

	// Parse request body with error handling
	req := &validation.ValidateRequest{
		PhoneNumber:    normalizePhoneNumber(c.FormValue("phone_number")),
		AccountID:      c.FormValue("account_id"),
		CustomerName:   c.FormValue("customer_name"),
		EventName:      c.FormValue("event_name"),
		WhatsappNumber: normalizePhoneNumber(c.FormValue("whatsapp_number")),
		Class:          c.FormValue("class"),
		Email:          c.FormValue("email"),
		SchoolType:     c.FormValue("school_type"),
	}

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

	_, errors := validation.ValidatePhoneNumber(req.PhoneNumber)
	fmt.Println("Phone Number:", req.PhoneNumber)
	fmt.Println("Errors:", errors)
	if errors != nil {
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	file, err := c.FormFile("registration_evidence")
	if err != nil {
		errors := map[string]string{
			"registration_evidence": "Bukti pendaftaran tidak boleh kosong",
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Validasi ekstensi file
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[ext] {
		errors := map[string]string{
			"registration_evidence": "File harus berupa gambar dengan ekstensi jpg, jpeg, atau png",
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// (Opsional) Validasi ukuran file, misal maksimum 5MB
	if file.Size > 5*1024*1024 {
		errors := map[string]string{
			"registration_evidence": "Ukuran file maksimum adalah 5MB",
		}
		helper.LogError(c, fmt.Sprintf("File size exceeds limit: %d bytes", file.Size))
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Simpan file
	uploadDir := "./public/uploads/registration-dealing"

	// Pastikan folder tersedia
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			helper.LogError(c, fmt.Sprintf("Failed to create directory: %v", err))
			response := helper.APIResponse("Gagal membuat direktori penyimpanan", fiber.StatusInternalServerError, "error", err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	// Generate filename unik berbasis timestamp
	timestamp := time.Now().UnixNano()
	ext = filepath.Ext(file.Filename) // ambil ekstensi asli
	filename := fmt.Sprintf("%d%s", timestamp, ext)
	filePath := filepath.Join(uploadDir, filename)

	// Simpan file
	if err := c.SaveFile(file, filePath); err != nil {
		helper.LogError(c, fmt.Sprintf("Failed to save file: %v", err))
		response := helper.APIResponse("Failed to save file", fiber.StatusInternalServerError, "error", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	req.RegistrationEvidence = &filePath

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
	var registrationDealing *response.RegistrationDealingResponse
	var serviceErr error

	select {
	case <-ctx.Done():
		response := helper.APIResponse("Request timeout during registration dealing creation", fiber.StatusRequestTimeout, "error", nil)
		return c.Status(fiber.StatusRequestTimeout).JSON(response)
	default:
		registrationDealing, serviceErr = h.service.CreateRegistrationDealing(reqMap, userID)
		if serviceErr != nil {
			helper.LogError(c, fmt.Sprintf("Failed to create registration dealing: %v", serviceErr))
			response := helper.APIResponse("Failed to create registration dealing: "+serviceErr.Error(), fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	// Return success response
	response := helper.APIResponse("Create Registration Dealing Succsesfully", fiber.StatusOK, "success", registrationDealing)
	return c.Status(fiber.StatusOK).JSON(response)
}

func normalizePhoneNumber(input string) string {
	input = strings.TrimSpace(input)
	input = strings.TrimPrefix(input, "+")
	if strings.HasPrefix(input, "62") {
		return input
	}
	if strings.HasPrefix(input, "0") {
		return "62" + input[1:]
	}
	return input
}
