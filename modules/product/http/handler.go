package http

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"byu-crm-service/helper"
	"byu-crm-service/models"
	eligibilityService "byu-crm-service/modules/eligibility/service"
	"byu-crm-service/modules/product/service"
	"byu-crm-service/modules/product/validation"

	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	service            service.ProductService
	eligibilityService eligibilityService.EligibilityService
}

func NewProductHandler(service service.ProductService, eligibilityService eligibilityService.EligibilityService) *ProductHandler {
	return &ProductHandler{service: service, eligibilityService: eligibilityService}
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
	accountID, _ := strconv.Atoi(c.Query("account_id", "0"))

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

func (h *ProductHandler) GetProductById(c *fiber.Ctx) error {
	idParam := c.Params("id")

	// Convert to int
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response := helper.APIResponse("Invalid ID format", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	product, err := h.service.FindByProductID(uint(id))
	if err != nil {
		response := helper.APIResponse("Product not found", fiber.StatusNotFound, "error", nil)
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	responseData := map[string]interface{}{
		"product": product,
	}

	response := helper.APIResponse("Success get product", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	// Add a timeout context to prevent long-running operations
	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	// Use a recovery function to catch any panics
	defer func() {
		if r := recover(); r != nil {
			helper.LogError(c, fmt.Sprintf("Panic in Create Product: %v", r))
			response := helper.APIResponse("Internal server error", fiber.StatusInternalServerError, "error", r)
			c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}()

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

	if req.ProductCategory == "VOUCHER" || req.ProductCategory == "PERDANA" || req.ProductCategory == "RENEWAL" {
		errors := validation.ValidateVoucherPerdana(req)
		if errors != nil {
			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	} else if req.ProductCategory == "SOLUSI" || req.ProductCategory == "LBO" {

		errors := validation.ValidateSolutionLbo(req)
		if errors != nil {
			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	} else if req.ProductCategory == "HOUSEHOLD" {

		errors := validation.ValidateHousehold(req)
		if errors != nil {
			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	} else {
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// if req.ProductCategory == "VOUCHER" || req.ProductCategory == "PERDANA" || req.ProductCategory == "RENEWAL" || req.ProductCategory == "HOUSEHOLD" {
	// 	*req.AdditionalFile = ""
	// 	keyVisualData := *req.KeyVisual

	// 	// Decode base64
	// 	decoded, mimeType, err := helper.DecodeBase64Image(keyVisualData)
	// 	if err != nil {
	// 		errors := map[string]string{
	// 			"key_visual": "Gambar tidak valid: " + err.Error(),
	// 		}
	// 		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
	// 		return c.Status(fiber.StatusBadRequest).JSON(response)
	// 	}

	// 	// Validasi ukuran maksimum 5MB
	// 	if len(decoded) > 5*1024*1024 {
	// 		errors := map[string]string{
	// 			"key_visual": "Ukuran file maksimum adalah 5MB",
	// 		}
	// 		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
	// 		return c.Status(fiber.StatusBadRequest).JSON(response)
	// 	}

	// 	// Validasi ekstensi
	// 	allowedMimes := map[string]string{
	// 		"image/jpeg": ".jpg",
	// 		"image/png":  ".png",
	// 	}
	// 	ext, ok := allowedMimes[mimeType]
	// 	if !ok {
	// 		errors := map[string]string{
	// 			"key_visual": "Format gambar tidak didukung",
	// 		}
	// 		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
	// 		return c.Status(fiber.StatusBadRequest).JSON(response)
	// 	}

	// 	// Simpan ke file
	// 	uploadDir := "./public/uploads/product"
	// 	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
	// 		response := helper.APIResponse("Failed to create directory", fiber.StatusInternalServerError, "error", err.Error())
	// 		return c.Status(fiber.StatusInternalServerError).JSON(response)
	// 	}

	// 	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	// 	filePath := filepath.Join(uploadDir, filename)

	// 	if err := os.WriteFile(filePath, decoded, 0644); err != nil {
	// 		response := helper.APIResponse("Failed to save file", fiber.StatusInternalServerError, "error", err.Error())
	// 		return c.Status(fiber.StatusInternalServerError).JSON(response)
	// 	}

	// 	req.KeyVisual = &filePath
	// } else if req.ProductCategory == "SOLUSI" || req.ProductCategory == "LBO" {
	// 	*req.KeyVisual = ""
	// 	AdditionalFileData := *req.AdditionalFile

	// 	// Decode base64 umum (tidak terbatas gambar)
	// 	decoded, mimeType, err := helper.DecodeBase64File(AdditionalFileData)
	// 	if err != nil {
	// 		errors := map[string]string{
	// 			"additional_file": "File tidak valid: " + err.Error(),
	// 		}
	// 		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
	// 		return c.Status(fiber.StatusBadRequest).JSON(response)
	// 	}

	// 	// Validasi ukuran maksimum 5MB
	// 	if len(decoded) > 5*1024*1024 {
	// 		errors := map[string]string{
	// 			"additional_file": "Ukuran file maksimum adalah 5MB",
	// 		}
	// 		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
	// 		return c.Status(fiber.StatusBadRequest).JSON(response)
	// 	}

	// 	// Ekstensi berdasarkan MIME (bisa ditambah sesuai kebutuhan)
	// 	mimeExtensions := map[string]string{
	// 		"image/jpeg":                    ".jpg",
	// 		"image/png":                     ".png",
	// 		"application/pdf":               ".pdf",
	// 		"application/zip":               ".zip",
	// 		"application/msword":            ".doc",
	// 		"application/vnd.ms-excel":      ".xls",
	// 		"application/vnd.ms-powerpoint": ".ppt",
	// 		"application/vnd.openxmlformats-officedocument.wordprocessingml.document":   ".docx",
	// 		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         ".xlsx",
	// 		"application/vnd.openxmlformats-officedocument.presentationml.presentation": ".pptx",
	// 		"text/plain":       ".txt",
	// 		"application/json": ".json",
	// 	}

	// 	// Dapatkan ekstensi
	// 	ext, ok := mimeExtensions[mimeType]
	// 	if !ok {
	// 		ext = ".bin" // default fallback
	// 	}

	// 	// Simpan ke file
	// 	uploadDir := "./public/uploads/product"
	// 	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
	// 		response := helper.APIResponse("Failed to create directory", fiber.StatusInternalServerError, "error", err.Error())
	// 		return c.Status(fiber.StatusInternalServerError).JSON(response)
	// 	}

	// 	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	// 	filePath := filepath.Join(uploadDir, filename)

	// 	if err := os.WriteFile(filePath, decoded, 0644); err != nil {
	// 		response := helper.APIResponse("Failed to save file", fiber.StatusInternalServerError, "error", err.Error())
	// 		return c.Status(fiber.StatusInternalServerError).JSON(response)
	// 	}

	// 	req.AdditionalFile = &filePath

	// }

	testValue := "testing"
	req.AdditionalFile = &testValue
	req.KeyVisual = &testValue

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
	var product *models.Product
	var serviceErr error

	select {
	case <-ctx.Done():
		response := helper.APIResponse("Request timeout during product creation", fiber.StatusRequestTimeout, "error", nil)
		return c.Status(fiber.StatusRequestTimeout).JSON(response)
	default:
		product, serviceErr = h.service.CreateProduct(reqMap)
		if serviceErr != nil {
			helper.LogError(c, fmt.Sprintf("Failed to create product: %v", serviceErr))
			response := helper.APIResponse("Failed to create product: "+serviceErr.Error(), fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	_ = h.eligibilityService.CreateEligibility("App\\Models\\Product", product.ID, req.EligibilityCategory, req.EligibilityType, req.EligibilityLocation)

	// Return success response
	response := helper.APIResponse("Create Product Succsesfully", fiber.StatusOK, "success", product)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	// Add a timeout context to prevent long-running operations
	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

	// Use a recovery function to catch any panics
	defer func() {
		if r := recover(); r != nil {
			helper.LogError(c, fmt.Sprintf("Panic in Update Product: %v", r))
			response := helper.APIResponse("Internal server error", fiber.StatusInternalServerError, "error", r)
			c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}()

	// Get and validate product ID
	productIDStr := c.Params("id")
	if productIDStr == "" {
		response := helper.APIResponse("Product ID is required", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		response := helper.APIResponse("Invalid Product ID", fiber.StatusBadRequest, "error", nil)
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

	if req.ProductCategory == "VOUCHER" || req.ProductCategory == "PERDANA" || req.ProductCategory == "RENEWAL" {
		errors := validation.ValidateVoucherPerdana(req)
		if errors != nil {
			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	} else if req.ProductCategory == "SOLUSI" || req.ProductCategory == "LBO" {

		errors := validation.ValidateSolutionLbo(req)
		if errors != nil {
			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	} else if req.ProductCategory == "HOUSEHOLD" {

		errors := validation.ValidateHousehold(req)
		if errors != nil {
			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
	} else {
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// if req.ProductCategory == "VOUCHER" || req.ProductCategory == "PERDANA" || req.ProductCategory == "RENEWAL" || req.ProductCategory == "HOUSEHOLD" {
	// 	*req.AdditionalFile = ""
	// 	if *req.KeyVisual != "" {
	// 		// helper.DeleteFile()
	// 		keyVisualData := *req.KeyVisual

	// 		// Decode base64
	// 		decoded, mimeType, err := helper.DecodeBase64Image(keyVisualData)
	// 		if err != nil {
	// 			errors := map[string]string{
	// 				"key_visual": "Gambar tidak valid: " + err.Error(),
	// 			}
	// 			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
	// 			return c.Status(fiber.StatusBadRequest).JSON(response)
	// 		}

	// 		// Validasi ukuran maksimum 5MB
	// 		if len(decoded) > 5*1024*1024 {
	// 			errors := map[string]string{
	// 				"key_visual": "Ukuran file maksimum adalah 5MB",
	// 			}
	// 			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
	// 			return c.Status(fiber.StatusBadRequest).JSON(response)
	// 		}

	// 		// Validasi ekstensi
	// 		allowedMimes := map[string]string{
	// 			"image/jpeg": ".jpg",
	// 			"image/png":  ".png",
	// 		}
	// 		ext, ok := allowedMimes[mimeType]
	// 		if !ok {
	// 			errors := map[string]string{
	// 				"key_visual": "Format gambar tidak didukung",
	// 			}
	// 			response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
	// 			return c.Status(fiber.StatusBadRequest).JSON(response)
	// 		}

	// 		// Simpan ke file
	// 		uploadDir := "./public/uploads/product"
	// 		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
	// 			response := helper.APIResponse("Failed to create directory", fiber.StatusInternalServerError, "error", err.Error())
	// 			return c.Status(fiber.StatusInternalServerError).JSON(response)
	// 		}

	// 		filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	// 		filePath := filepath.Join(uploadDir, filename)

	// 		if err := os.WriteFile(filePath, decoded, 0644); err != nil {
	// 			response := helper.APIResponse("Failed to save file", fiber.StatusInternalServerError, "error", err.Error())
	// 			return c.Status(fiber.StatusInternalServerError).JSON(response)
	// 		}

	// 		req.KeyVisual = &filePath
	// 	} else {
	// 		*req.KeyVisual = ""
	// 	}
	// } else if req.ProductCategory == "SOLUSI" || req.ProductCategory == "LBO" {
	// 	*req.KeyVisual = ""
	// 	AdditionalFileData := *req.AdditionalFile

	// 	// Decode base64 umum (tidak terbatas gambar)
	// 	decoded, mimeType, err := helper.DecodeBase64File(AdditionalFileData)
	// 	if err != nil {
	// 		errors := map[string]string{
	// 			"additional_file": "File tidak valid: " + err.Error(),
	// 		}
	// 		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
	// 		return c.Status(fiber.StatusBadRequest).JSON(response)
	// 	}

	// 	// Validasi ukuran maksimum 5MB
	// 	if len(decoded) > 5*1024*1024 {
	// 		errors := map[string]string{
	// 			"additional_file": "Ukuran file maksimum adalah 5MB",
	// 		}
	// 		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
	// 		return c.Status(fiber.StatusBadRequest).JSON(response)
	// 	}

	// 	// Ekstensi berdasarkan MIME (bisa ditambah sesuai kebutuhan)
	// 	mimeExtensions := map[string]string{
	// 		"image/jpeg":                    ".jpg",
	// 		"image/png":                     ".png",
	// 		"application/pdf":               ".pdf",
	// 		"application/zip":               ".zip",
	// 		"application/msword":            ".doc",
	// 		"application/vnd.ms-excel":      ".xls",
	// 		"application/vnd.ms-powerpoint": ".ppt",
	// 		"application/vnd.openxmlformats-officedocument.wordprocessingml.document":   ".docx",
	// 		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         ".xlsx",
	// 		"application/vnd.openxmlformats-officedocument.presentationml.presentation": ".pptx",
	// 		"text/plain":       ".txt",
	// 		"application/json": ".json",
	// 	}

	// 	// Dapatkan ekstensi
	// 	ext, ok := mimeExtensions[mimeType]
	// 	if !ok {
	// 		ext = ".bin" // default fallback
	// 	}

	// 	// Simpan ke file
	// 	uploadDir := "./public/uploads/product"
	// 	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
	// 		response := helper.APIResponse("Failed to create directory", fiber.StatusInternalServerError, "error", err.Error())
	// 		return c.Status(fiber.StatusInternalServerError).JSON(response)
	// 	}

	// 	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	// 	filePath := filepath.Join(uploadDir, filename)

	// 	if err := os.WriteFile(filePath, decoded, 0644); err != nil {
	// 		response := helper.APIResponse("Failed to save file", fiber.StatusInternalServerError, "error", err.Error())
	// 		return c.Status(fiber.StatusInternalServerError).JSON(response)
	// 	}

	// 	req.AdditionalFile = &filePath

	// }

	testValue := "testing"
	req.AdditionalFile = &testValue
	req.KeyVisual = &testValue

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
	var product *models.Product
	var serviceErr error

	select {
	case <-ctx.Done():
		response := helper.APIResponse("Request timeout during product update", fiber.StatusRequestTimeout, "error", nil)
		return c.Status(fiber.StatusRequestTimeout).JSON(response)
	default:
		product, serviceErr = h.service.UpdateProduct(reqMap, productID)
		if serviceErr != nil {
			helper.LogError(c, fmt.Sprintf("Failed to update product: %v", serviceErr))
			response := helper.APIResponse("Failed to update product: "+serviceErr.Error(), fiber.StatusInternalServerError, "error", nil)
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	_ = h.eligibilityService.CreateEligibility("App\\Models\\Product", product.ID, req.EligibilityCategory, req.EligibilityType, req.EligibilityLocation)

	// Return success response
	response := helper.APIResponse("Update Product Succsesfully", fiber.StatusOK, "success", product)
	return c.Status(fiber.StatusOK).JSON(response)
}
