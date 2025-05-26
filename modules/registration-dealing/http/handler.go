package http

import (
	"log"
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

func (h *RegistrationDealingHandler) GetAllRegistrationDealingsGrouped(c *fiber.Ctx) error {
	// Default query params
	filters := map[string]string{
		"search":     c.Query("search", ""),
		"order_by":   c.Query("order_by", "id"),
		"order":      c.Query("order", "ASC"),
		"start_date": c.Query("start_date", ""),
		"end_date":   c.Query("end_date", ""),
	}

	// Parse integer and boolean values
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	paginate, _ := strconv.ParseBool(c.Query("paginate", "true"))
	page, _ := strconv.Atoi(c.Query("page", "1"))

	// Call service with filters
	registrationDealings, total, err := h.service.GetAllRegistrationDealingGrouped(limit, paginate, page, filters)
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

	responseData := map[string]interface{}{
		"registration_dealing": registrationDealing,
	}

	response := helper.APIResponse("Success get registration dealing", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *RegistrationDealingHandler) CreateRegistrationDealing(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	defer func() {
		if r := recover(); r != nil {
			log.Printf(fmt.Sprintf("Panic in Create Registration Dealing: %v", r))
			response := helper.APIResponse("Internal server error", fiber.StatusInternalServerError, "error", r)
			c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}()

	var userID *int
	if val, ok := c.Locals("user_id").(int); ok {
		userID = &val
	}

	req := new(validation.ValidateRequest)

	if err := c.BodyParser(req); err != nil {
		msg := "Invalid request format: " + err.Error()
		if err.Error() == "unexpected EOF" {
			msg = "Invalid request: Unexpected end of JSON input"
		}
		response := helper.APIResponse(msg, fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	req.PhoneNumber = normalizePhoneNumber(req.PhoneNumber)
	req.WhatsappNumber = normalizePhoneNumber(req.WhatsappNumber)

	// Validasi data umum
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

	// Validasi nomor HP
	if _, errors := validation.ValidatePhoneNumber(req.PhoneNumber); errors != nil {
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Validasi dan simpan base64 image
	if req.RegistrationEvidence == nil || *req.RegistrationEvidence == "" {
		errors := map[string]string{
			"registration_evidence": "Bukti pendaftaran tidak boleh kosong",
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	imageData := *req.RegistrationEvidence

	// Decode base64
	decoded, mimeType, err := helper.DecodeBase64Image(imageData)
	if err != nil {
		errors := map[string]string{
			"registration_evidence": "Gambar tidak valid: " + err.Error(),
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Validasi ukuran maksimum 5MB
	if len(decoded) > 5*1024*1024 {
		errors := map[string]string{
			"registration_evidence": "Ukuran file maksimum adalah 5MB",
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Validasi ekstensi
	allowedMimes := map[string]string{
		"image/jpeg": ".jpg",
		"image/png":  ".png",
	}
	ext, ok := allowedMimes[mimeType]
	if !ok {
		errors := map[string]string{
			"registration_evidence": "Format gambar tidak didukung",
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Simpan ke file
	uploadDir := "./public/uploads/registration-dealing"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		response := helper.APIResponse("Failed to create directory", fiber.StatusInternalServerError, "error", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(uploadDir, filename)

	if err := os.WriteFile(filePath, decoded, 0644); err != nil {
		response := helper.APIResponse("Failed to save file", fiber.StatusInternalServerError, "error", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	req.RegistrationEvidence = &filePath

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
			log.Printf(fmt.Sprintf("Failed to marshal request: %v", marshalErr))
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
			log.Printf(fmt.Sprintf("Failed to unmarshal request: %v", unmarshalErr))
			response := helper.APIResponse("Failed to process request data", fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	var registrationDealing *response.RegistrationDealingResponse
	var serviceErr error
	select {
	case <-ctx.Done():
		response := helper.APIResponse("Request timeout during registration dealing creation", fiber.StatusRequestTimeout, "error", nil)
		return c.Status(fiber.StatusRequestTimeout).JSON(response)
	default:
		registrationDealing, serviceErr = h.service.CreateRegistrationDealing(reqMap, userID)
		if serviceErr != nil {
			log.Printf(fmt.Sprintf("Failed to create account: %v", serviceErr))
			response := helper.APIResponse("Failed to create account: "+serviceErr.Error(), fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	return c.Status(fiber.StatusOK).JSON(helper.APIResponse("Create Registration Dealing Successfully", fiber.StatusOK, "success", registrationDealing))
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
