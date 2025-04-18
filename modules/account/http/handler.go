package http

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"byu-crm-service/helper"
	"byu-crm-service/modules/account/service"
	"byu-crm-service/modules/account/validation"

	accountFacultyService "byu-crm-service/modules/account-faculty/service"
	accountMemberService "byu-crm-service/modules/account-member/service"
	accountScheduleService "byu-crm-service/modules/account-schedule/service"
	accountTypeCampusDetailService "byu-crm-service/modules/account-type-campus-detail/service"
	accountTypeCommunityDetailService "byu-crm-service/modules/account-type-community-detail/service"
	accountTypeSchoolDetailService "byu-crm-service/modules/account-type-school-detail/service"
	contactAccountService "byu-crm-service/modules/contact-account/service"
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
	accountTypeCommunityDetailService accountTypeCommunityDetailService.AccountTypeCommunityDetailService) *AccountHandler {

	return &AccountHandler{
		service:                           service,
		contactAccountService:             contactAccountService,
		socialMediaService:                socialMediaService,
		accountTypeSchoolDetailService:    accountTypeSchoolDetailService,
		accountFacultyService:             accountFacultyService,
		accountMemberService:              accountMemberService,
		accountScheduleService:            accountScheduleService,
		accountTypeCampusDetailService:    accountTypeCampusDetailService,
		accountTypeCommunityDetailService: accountTypeCommunityDetailService}
}

func (h *AccountHandler) GetAllAccounts(c *fiber.Ctx) error {
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
	// Parse the multipart form data
	form, err := c.MultipartForm()
	if err != nil {
		response := helper.APIResponse("Failed to parse form", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Initialize a map to store the form data
	requestBody := make(map[string]interface{})

	// Loop through the form fields and store them in the requestBody map
	for key, values := range form.Value {
		// If the field has only one value, take the first element to avoid array
		if len(values) == 1 {
			requestBody[key] = values[0]
		} else {
			// Otherwise, store the array as is
			requestBody[key] = values
		}
	}

	// // Handle file upload for account image
	// accountImage, err := c.FormFile("account_image")
	// if err == nil && accountImage != nil {
	// 	// Allowed image formats
	// 	allowedFormats := []string{".jpg", ".jpeg", ".png", ".gif"}

	// 	// Save the uploaded image
	// 	accountImagePath, err := saveFileToLocal(accountImage, "uploads/account_images", allowedFormats)
	// 	if err != nil {
	// 		response := helper.APIResponse("Failed to save image", fiber.StatusInternalServerError, "error", nil)
	// 		return c.Status(fiber.StatusInternalServerError).JSON(response)
	// 	}

	// 	// Add the image path to the request body
	// 	requestBody["account_image"] = *accountImagePath
	// }

	// // Validate and parse the form data into CreateRequest
	// req := new(validation.CreateRequest)
	// // We need to manually populate req with data from requestBody
	// if err := mapstructure.Decode(requestBody, req); err != nil {
	// 	response := helper.APIResponse("Invalid request data", fiber.StatusBadRequest, "error", nil)
	// 	return c.Status(fiber.StatusBadRequest).JSON(response)
	// }

	// // Request Validation
	// errors := validation.ValidateCreate(req)
	// if errors != nil {
	// 	response := helper.APIResponse("Validation error", fiber.StatusBadRequest, "error", errors)
	// 	return c.Status(fiber.StatusBadRequest).JSON(response)
	// }

	// Extract user ID from query parameter
	userID := 1
	// if err != nil || userID == 0 {
	// response := helper.APIResponse("Invalid user ID", fiber.StatusBadRequest, "error", requestBody)
	// return c.Status(fiber.StatusBadRequest).JSON(response)
	// }

	// Create Account
	account, err := h.service.CreateAccount(requestBody, userID)
	if err != nil {
		response := helper.APIResponse("Failed to create account", fiber.StatusInternalServerError, "error", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	if len(account) > 0 {
		_, _ = h.contactAccountService.InsertContactAccount(requestBody, account[0].ID)
		_, _ = h.socialMediaService.InsertSocialMedia(requestBody, "App\\Models\\Account", account[0].ID)
		if category, exists := requestBody["account_category"]; exists {
			if categoryStr, ok := category.(string); ok {
				switch categoryStr {
				case "SEKOLAH":
					_, _ = h.accountTypeSchoolDetailService.Insert(requestBody, account[0].ID)
				case "KAMPUS":
					_, _ = h.accountFacultyService.Insert(requestBody, account[0].ID)
					_, _ = h.accountMemberService.Insert(requestBody, "App\\Models\\Account", account[0].ID, "year", "amount")
					_, _ = h.accountMemberService.Insert(requestBody, "App\\Models\\AccountLecture", account[0].ID, "year_lecture", "amount_lecture")
					_, _ = h.accountScheduleService.Insert(requestBody, "App\\Models\\Account", account[0].ID)
					_, _ = h.accountTypeCampusDetailService.Insert(requestBody, account[0].ID)
				case "KOMUNITAS":
					_, _ = h.accountTypeCommunityDetailService.Insert(requestBody, account[0].ID)
					_, _ = h.accountScheduleService.Insert(requestBody, "App\\Models\\Account", account[0].ID)
				}
			}
		}
	}

	// Return success response
	response := helper.APIResponse("Account successfully created", fiber.StatusOK, "success", account)
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *AccountHandler) UpdateAccount(c *fiber.Ctx) error {
	territoryID := c.Locals("territory_id").(int)
	userRole := c.Locals("user_role").(string)
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

	// Parse the multipart form data
	form, err := c.MultipartForm()
	if err != nil {
		response := helper.APIResponse("Failed to parse form", fiber.StatusBadRequest, "error", nil)
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Initialize a map to store the form data
	requestBody := make(map[string]interface{})

	// Loop through the form fields and store them in the requestBody map
	for key, values := range form.Value {
		if len(values) == 1 {
			requestBody[key] = values[0]
		} else {
			requestBody[key] = values
		}
	}

	// Handle "byod" field (convert from string to uint)
	if byodStr, exists := requestBody["byod"].(string); exists {
		byodValue, err := strconv.ParseUint(byodStr, 10, 32)
		if err != nil {
			byodValue = 0 // Default ke 0 jika gagal parsing
		}
		requestBody["byod"] = uint(byodValue)
	}

	// Ambil user ID dari sesi atau token (sementara default ke 1)
	userID := 1

	// Panggil service untuk update account
	account, err := h.service.UpdateAccount(requestBody, accountID, userRole, territoryID, userID)
	if err != nil {
		response := helper.APIResponse("Failed to update account", fiber.StatusInternalServerError, "error", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	if len(account) > 0 {
		_, _ = h.contactAccountService.InsertContactAccount(requestBody, account[0].ID)
		_, _ = h.socialMediaService.InsertSocialMedia(requestBody, "App\\Models\\Account", account[0].ID)
		if category, exists := requestBody["account_category"]; exists {
			if categoryStr, ok := category.(string); ok {
				switch categoryStr {
				case "SEKOLAH":
					_, _ = h.accountTypeSchoolDetailService.Insert(requestBody, account[0].ID)
				case "KAMPUS":
					_, _ = h.accountFacultyService.Insert(requestBody, account[0].ID)
					_, _ = h.accountMemberService.Insert(requestBody, "App\\Models\\Account", account[0].ID, "year", "amount")
					_, _ = h.accountMemberService.Insert(requestBody, "App\\Models\\AccountLecture", account[0].ID, "year_lecture", "amount_lecture")
					_, _ = h.accountScheduleService.Insert(requestBody, "App\\Models\\Account", account[0].ID)
					_, _ = h.accountTypeCampusDetailService.Insert(requestBody, account[0].ID)
				case "KOMUNITAS":
					_, _ = h.accountTypeCommunityDetailService.Insert(requestBody, account[0].ID)
					_, _ = h.accountScheduleService.Insert(requestBody, "App\\Models\\Account", account[0].ID)
				}
			}
		}
	}

	// Return success response
	response := helper.APIResponse("Account successfully updated", fiber.StatusOK, "success", account)
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
