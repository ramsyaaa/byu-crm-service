package helper

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"byu-crm-service/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// DatabaseLogger creates a middleware that logs API requests to the database
func DatabaseLogger(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Skip logging for internal admin API calls
		if strings.HasPrefix(c.OriginalURL(), "/admin/") {
			return c.Next()
		}

		// Store the start time
		startTime := time.Now()

		// Capture request body
		var requestBody []byte
		if c.Body() != nil {
			requestBody = make([]byte, len(c.Body()))
			copy(requestBody, c.Body())
		}

		// Process the request
		err := c.Next()

		// Calculate response time
		responseTime := time.Since(startTime).Milliseconds()

		// Extract user email from JWT token if available
		var userEmail *string
		if authHeader := c.Get("Authorization"); authHeader != "" {
			if email := extractEmailFromJWT(authHeader); email != "" {
				userEmail = &email
			}
		}

		// Prepare request payload
		var requestPayload *string
		if len(requestBody) > 0 && len(requestBody) < 10000 { // Limit size to prevent huge payloads
			requestStr := string(requestBody)
			requestPayload = &requestStr
		}

		// Prepare response payload
		var responsePayload *string
		responseBody := c.Response().Body()
		if len(responseBody) > 0 && len(responseBody) < 10000 { // Limit size to prevent huge payloads
			responseStr := string(responseBody)
			responsePayload = &responseStr
		}

		// Prepare error message
		var errorMessage *string
		if err != nil {
			errMsg := err.Error()
			errorMessage = &errMsg
		} else if c.Response().StatusCode() >= 400 {
			// Capture error message from response body for HTTP error responses
			responseBody := c.Response().Body()
			if len(responseBody) > 0 && len(responseBody) < 5000 { // Limit size
				// Try to extract error message from JSON response
				if errorMsg := extractErrorFromResponse(responseBody); errorMsg != "" {
					errorMessage = &errorMsg
				}
			}
		}

		// Create log entry
		logEntry := models.ApiLog{
			AccessedAt:      startTime,
			Endpoint:        c.OriginalURL(),
			Method:          c.Method(),
			StatusCode:      c.Response().StatusCode(),
			ResponseTimeMs:  responseTime,
			RequestPayload:  requestPayload,
			ResponsePayload: responsePayload,
			ErrorMessage:    errorMessage,
			AuthUserEmail:   userEmail,
			IPAddress:       c.IP(),
			UserAgent:       c.Get("User-Agent"),
		}

		// Log to database asynchronously
		go func() {
			if dbErr := db.Create(&logEntry).Error; dbErr != nil {
				// Fallback to console logging if database fails
				log.Printf("Failed to log to database: %v. Log entry: %+v", dbErr, logEntry)
			}
		}()

		return err
	}
}

// extractErrorFromResponse tries to extract error message from JSON response body
func extractErrorFromResponse(responseBody []byte) string {
	// Try to parse as JSON and extract common error fields
	var jsonResponse map[string]interface{}
	if err := json.Unmarshal(responseBody, &jsonResponse); err != nil {
		// If not JSON, return the raw response (truncated)
		response := string(responseBody)
		if len(response) > 200 {
			response = response[:200] + "..."
		}
		return response
	}

	// Try common error message fields
	errorFields := []string{"message", "error", "error_message", "msg", "detail", "description"}
	for _, field := range errorFields {
		if value, exists := jsonResponse[field]; exists {
			if str, ok := value.(string); ok && str != "" {
				return str
			}
		}
	}

	// If no specific error field found, try to get a general message
	if status, exists := jsonResponse["status"]; exists {
		if str, ok := status.(string); ok && str == "error" {
			// Look for any string value that might be an error
			for key, value := range jsonResponse {
				if key != "status" {
					if str, ok := value.(string); ok && str != "" {
						return str
					}
				}
			}
		}
	}

	return ""
}

// extractEmailFromJWT extracts email from JWT token in Authorization header
func extractEmailFromJWT(authHeader string) string {
	// Remove "Bearer " prefix
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return "" // No Bearer prefix found
	}

	// Parse token without verification (we just want to extract claims)
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return ""
	}

	// Extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if email, exists := claims["email"]; exists {
			if emailStr, ok := email.(string); ok {
				return emailStr
			}
		}
	}

	return ""
}

// LogErrorToDatabase logs an error message to the database
func LogErrorToDatabase(db *gorm.DB, c *fiber.Ctx, errorMsg string) {
	// Skip logging for internal admin API calls
	if strings.HasPrefix(c.OriginalURL(), "/admin/") {
		return
	}

	// Extract user email from JWT token if available
	var userEmail *string
	if authHeader := c.Get("Authorization"); authHeader != "" {
		if email := extractEmailFromJWT(authHeader); email != "" {
			userEmail = &email
		}
	}

	// Create error log entry
	logEntry := models.ApiLog{
		AccessedAt:      time.Now(),
		Endpoint:        c.OriginalURL(),
		Method:          c.Method(),
		StatusCode:      500, // Internal Server Error
		ResponseTimeMs:  0,   // Unknown for errors
		RequestPayload:  nil,
		ResponsePayload: nil,
		ErrorMessage:    &errorMsg,
		AuthUserEmail:   userEmail,
		IPAddress:       c.IP(),
		UserAgent:       c.Get("User-Agent"),
	}

	// Log to database asynchronously
	go func() {
		if dbErr := db.Create(&logEntry).Error; dbErr != nil {
			// Fallback to console logging if database fails
			log.Printf("Failed to log error to database: %v. Error: %s", dbErr, errorMsg)
		}
	}()
}

// LogPanicToDatabase logs a panic message to the database
func LogPanicToDatabase(db *gorm.DB, c *fiber.Ctx, panicMsg string) {
	// Extract user email from JWT token if available
	var userEmail *string
	if c != nil {
		if authHeader := c.Get("Authorization"); authHeader != "" {
			if email := extractEmailFromJWT(authHeader); email != "" {
				userEmail = &email
			}
		}
	}

	endpoint := "-"
	method := "-"
	ipAddress := "-"
	userAgent := "-"

	if c != nil {
		endpoint = c.OriginalURL()
		method = c.Method()
		ipAddress = c.IP()
		userAgent = c.Get("User-Agent")
	}

	// Create panic log entry
	logEntry := models.ApiLog{
		AccessedAt:      time.Now(),
		Endpoint:        endpoint,
		Method:          method,
		StatusCode:      500, // Internal Server Error
		ResponseTimeMs:  0,   // Unknown for panics
		RequestPayload:  nil,
		ResponsePayload: nil,
		ErrorMessage:    &panicMsg,
		AuthUserEmail:   userEmail,
		IPAddress:       ipAddress,
		UserAgent:       userAgent,
	}

	// Log to database asynchronously
	go func() {
		if dbErr := db.Create(&logEntry).Error; dbErr != nil {
			// Fallback to console logging if database fails
			log.Printf("Failed to log panic to database: %v. Panic: %s", dbErr, panicMsg)
		}
	}()
}
