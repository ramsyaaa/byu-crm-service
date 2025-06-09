package http

import (
	"bytes"
	"byu-crm-service/helper"
	"byu-crm-service/modules/performance-skul-id/service"
	"byu-crm-service/modules/performance-skul-id/validation"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
)

type PerformanceSkulIdHandler struct {
	service service.PerformanceSkulIdService
}

func NewPerformanceSkulIdHandler(service service.PerformanceSkulIdService) *PerformanceSkulIdHandler {
	return &PerformanceSkulIdHandler{service: service}
}

func (h *PerformanceSkulIdHandler) Import(c *fiber.Ctx) error {
	// Validate the uploaded file
	if err := validation.ValidatePerformanceSkulIdRequest(c); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Save file temporarily
	file, err := c.FormFile("file_csv")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to retrieve file"})
	}

	tempPath := "./temp/" + file.Filename

	// Ensure the temp directory exists
	if _, err := os.Stat("./temp/"); os.IsNotExist(err) {
		if err := os.Mkdir("./temp/", os.ModePerm); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create temp directory"})
		}
	}

	// Save file
	if err := c.SaveFile(file, tempPath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save file"})
	}

	// Retrieve user_id from the form
	userID := c.FormValue("user_id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID is required"})
	}

	// Menghitung jumlah total baris untuk estimasi durasi
	totalRows, err := countCSVRows(tempPath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to count CSV rows"})
	}

	// Asumsi setiap baris membutuhkan 0.5 detik untuk diproses
	processingTimePerRow := 0.5
	estimatedDuration := time.Duration(float64(totalRows) * processingTimePerRow * float64(time.Second))

	// Respond immediately to the user
	go func() {
		defer os.Remove(tempPath) // Clean up the temporary file

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
		rows, err := reader.ReadAll()
		if err != nil {
			fmt.Println("Failed to read CSV:", err)
			return
		}

		for i, row := range rows {
			if i == 0 {
				continue // Skip header
			}
			if err := h.service.ProcessPerformanceSkulId(row); err != nil {
				fmt.Println("Error processing row:", err)
				continue // Lanjutkan proses meskipun ada error pada satu baris
			}
		}

		// Send notification
		notificationURL := os.Getenv("NOTIFICATION_URL") + "/api/notification/create"
		payload := map[string]interface{}{
			"model":    "App\\Models\\PerformanceSkulId",
			"model_id": 0, // Replace with actual model ID if needed
			"user_id":  userID,
			"data": map[string]string{
				"title":        "Import Performance SkulId",
				"description":  "Import Performance SkulId",
				"callback_url": "/performances-skulId",
			},
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			fmt.Println("Failed to marshal notification payload:", err)
			return
		}

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Post(notificationURL, "application/json", bytes.NewReader(payloadBytes))
		if err != nil {
			fmt.Println("Failed to send notification:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println("Notification API responded with status:", resp.StatusCode)
			return
		}

		var responseMap map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
			fmt.Println("Failed to decode response:", err)
			return
		}
		fmt.Println("Notification sent successfully:", responseMap["message"])
	}()

	return c.JSON(fiber.Map{
		"message":           "File upload successful, processing in background",
		"estimated_seconds": estimatedDuration.Seconds(),
	})
}

func (h *PerformanceSkulIdHandler) ImportByAccount(c *fiber.Ctx) error {
	// Validasi form
	if err := validation.ValidatePerformanceSkulIdRequest(c); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	userID := c.Locals("user_id").(int)
	territoryID := c.Locals("territory_id").(int)
	userRole := c.Locals("user_role").(string)
	idParam := c.Params("id")

	accountID, err := strconv.Atoi(idParam)
	if err != nil {
		response := helper.APIResponse("Invalid Account ID format", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	userType := c.FormValue("user_type")
	if strings.TrimSpace(userType) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user_type tidak boleh kosong"})
	}

	if userType != "Siswa" && userType != "Sekolah" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user_type harus 'Siswa' atau 'Sekolah'"})
	}

	// Ambil base64 dari form
	base64File := c.FormValue("file_import")
	if strings.TrimSpace(base64File) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File harus diunggah dalam bentuk base64"})
	}

	// Decode file base64
	decoded, mimeType, err := helper.DecodeBase64File(base64File)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File tidak valid: " + err.Error()})
	}

	// Validasi ukuran maksimal 5MB
	if len(decoded) > 5*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Ukuran file maksimum adalah 5MB"})
	}

	// Validasi MIME type
	var ext string
	switch mimeType {
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		ext = ".xlsx"
	case "application/vnd.ms-excel", "text/csv":
		ext = ".csv"
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File harus berformat CSV atau XLSX"})
	}

	// Simpan file sementara
	tempPath := fmt.Sprintf("./temp/upload_%d%s", time.Now().UnixNano(), ext)
	if _, err := os.Stat("./temp/"); os.IsNotExist(err) {
		_ = os.Mkdir("./temp/", os.ModePerm)
	}
	if err := os.WriteFile(tempPath, decoded, 0644); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan file sementara"})
	}

	// Hitung estimasi proses
	totalRows, err := countRows(tempPath, ext)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menghitung jumlah baris data"})
	}
	estimatedDuration := time.Duration(float64(totalRows)*0.5) * time.Second

	// Proses di background
	go func() {
		defer os.Remove(tempPath)

		var rows [][]string

		if ext == ".csv" {
			f, err := os.Open(tempPath)
			if err != nil {
				fmt.Println("Gagal membuka file CSV:", err)
				return
			}
			defer f.Close()

			reader := csv.NewReader(f)
			rows, err = reader.ReadAll()
			if err != nil {
				fmt.Println("Gagal membaca file CSV:", err)
				return
			}

		} else if ext == ".xlsx" {
			f, err := excelize.OpenFile(tempPath)
			if err != nil {
				fmt.Println("Gagal membuka file Excel:", err)
				return
			}
			defer f.Close()

			sheet := f.GetSheetName(0)
			rows, err = f.GetRows(sheet)
			if err != nil {
				fmt.Println("Gagal membaca baris Excel:", err)
				return
			}
		}

		for i, row := range rows {
			if i == 0 {
				continue // Lewati header
			}
			if err := h.service.ProcessPerformanceSkulIdByAccount(row, userID, accountID, userRole, territoryID, userType); err != nil {
				fmt.Println("Gagal memproses baris:", err)
				continue
			}
		}
	}()

	response := helper.APIResponse(
		"File berhasil diunggah dan sedang diproses di latar belakang",
		fiber.StatusOK,
		"success",
		map[string]interface{}{
			"estimated_seconds": estimatedDuration.Seconds(),
		},
	)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *PerformanceSkulIdHandler) CreatePerformanceSkulID(c *fiber.Ctx) error {
	req := new(validation.CreatePerformanceSkulIdRequest)
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

	idParam := c.Params("account_id")

	// Convert to int
	account_id, err := strconv.Atoi(idParam)
	if err != nil {
		response := helper.APIResponse("Invalid ID format", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	existingIdSkulId, _ := h.service.FindByIdSkulId(req.IDSkulId)
	if existingIdSkulId != nil {
		errors := map[string]string{
			"id_skulid": "ID SkulID sudah digunakan",
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	existingMSISDN, _ := h.service.FindBySerialNumberMsisdn(normalizePhoneNumber(req.MSISDN))
	if existingMSISDN != nil {
		errors := map[string]string{
			"msisdn": "MSISDN sudah digunakan",
		}
		response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Parse RegisteredDate from string to *time.Time
	var registeredDate *time.Time
	if req.RegisteredDate != "" {
		parsedDate, err := time.Parse("2006-01-02", req.RegisteredDate)
		if err != nil {
			response := helper.APIResponse("Invalid date format for RegisteredDate, expected YYYY-MM-DD", fiber.StatusBadRequest, "error", nil)
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
		registeredDate = &parsedDate
	}

	performance, err := h.service.CreatePerformanceSkulID(account_id, req.UserName, req.IDSkulId, normalizePhoneNumber(req.MSISDN), registeredDate, &req.Provider, &req.Batch, &req.UserType)
	if err != nil {
		response := helper.APIResponse(err.Error(), fiber.StatusUnauthorized, "error", nil)
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}

	// Response
	response := helper.APIResponse("Performance created successful", fiber.StatusOK, "success", performance)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *PerformanceSkulIdHandler) GetAllSkulIds(c *fiber.Ctx) error {
	// Default query params
	var account_id int
	paramAccountID := c.Query("account_id")
	if parsedID, err := strconv.Atoi(paramAccountID); err == nil {
		account_id = parsedID
	} else {
		response := helper.APIResponse("Invalid account_id", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	filters := map[string]string{
		"search":     c.Query("search", ""),
		"order_by":   c.Query("order_by", "id"),
		"order":      c.Query("order", "DESC"),
		"start_date": c.Query("start_date", ""),
		"end_date":   c.Query("end_date", ""),
		"user_type":  c.Query("user_type", ""),
	}

	// Parse integer and boolean values
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	paginate, _ := strconv.ParseBool(c.Query("paginate", "true"))
	page, _ := strconv.Atoi(c.Query("page", "1"))

	// Call service with filters
	performance_skulid, total, err := h.service.FindAll(limit, (page-1)*limit, filters, account_id, page, paginate)
	if err != nil {
		response := helper.APIResponse("Failed to fetch performance", fiber.StatusInternalServerError, "error", nil)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Return response
	responseData := map[string]interface{}{
		"performance_skulid": performance_skulid,
		"total":              total,
		"page":               page,
	}

	response := helper.APIResponse("Get Performance Skul ID Successfully", fiber.StatusOK, "success", responseData)
	return c.Status(fiber.StatusOK).JSON(response)
}

// countCSVRows menghitung jumlah total baris dalam file CSV
func countCSVRows(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	totalRows := 0

	// Baca setiap baris dan hitung jumlah totalnya
	for {
		_, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return 0, err
		}
		totalRows++
	}

	return totalRows, nil
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

func countRows(path, ext string) (int, error) {
	if ext == ".csv" {
		f, err := os.Open(path)
		if err != nil {
			return 0, err
		}
		defer f.Close()

		reader := csv.NewReader(f)
		rows, err := reader.ReadAll()
		if err != nil {
			return 0, err
		}
		return len(rows) - 1, nil // minus header

	} else if ext == ".xlsx" {
		f, err := excelize.OpenFile(path)
		if err != nil {
			return 0, err
		}
		defer f.Close()

		sheet := f.GetSheetName(0)
		rows, err := f.GetRows(sheet)
		if err != nil {
			return 0, err
		}
		return len(rows) - 1, nil
	}

	return 0, fmt.Errorf("unsupported file extension")
}
