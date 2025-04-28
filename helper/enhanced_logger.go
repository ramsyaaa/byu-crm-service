package helper

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// EnhancedLogToFile creates a logger middleware that includes the API response message
func LogToFile() fiber.Handler {
	// Determine the log file path based on current date
	filePath := "logs/" + time.Now().Format("20060102") + "-log.log"

	// Check if the log file for today already exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// If the file doesn't exist, create a new one
		if err := os.MkdirAll(filepath.Dir(filePath), 0666); err != nil {
			log.Fatalf("error creating directory: %v", err)
		}
		if _, err := os.Create(filePath); err != nil {
			log.Fatalf("error creating file: %v", err)
		}
	}

	// Open the file with read/write, create or append mode
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	// Create a middleware chain
	return func(c *fiber.Ctx) error {
		// Store the start time
		startTime := time.Now()
		c.Locals("start_time", startTime)

		// Process the request
		err := c.Next()

		// After the request is processed, log the details
		// Get the response body
		body := c.Response().Body()
		message := "-" // Default value if no message found

		if len(body) > 0 {
			// Try to parse the response as JSON
			var response struct {
				Meta struct {
					Message string `json:"message"`
				} `json:"meta"`
			}
			if err := json.Unmarshal(body, &response); err == nil {
				// If we have a meta message, use it
				if response.Meta.Message != "" {
					message = response.Meta.Message
				}
			}
		}

		// Create a new log entry with the response message
		statusCode := strconv.Itoa(c.Response().Header.StatusCode())
		latency := time.Since(startTime).String()

		logEntry := time.Now().Format(time.RFC3339) +
			" | " + statusCode +
			" | " + latency +
			" | " + c.IP() +
			" | " + c.Method() +
			" | " + c.Path() +
			" | " + message + "\n"

		// Append the log entry to the file
		if _, err := file.WriteString(logEntry); err != nil {
			log.Printf("Error writing to log file: %v", err)
		}

		return err
	}
}
