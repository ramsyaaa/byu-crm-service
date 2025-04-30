package helper

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Global file handle for logging
var (
	logFile     *os.File
	logFilePath string
	logMutex    sync.Mutex
)

// getLogFile returns the current log file handle, creating it if necessary
func getLogFile() *os.File {
	logMutex.Lock()
	defer logMutex.Unlock()

	// Check if we need to create or rotate the log file
	currentFilePath := "logs/" + time.Now().Format("20060102") + "-log.log"

	// If the file path has changed or the file handle is nil, create a new file
	if currentFilePath != logFilePath || logFile == nil {
		// Close the previous file if it exists
		if logFile != nil {
			logFile.Close()
		}

		// Check if the log file for today already exists
		if _, err := os.Stat(currentFilePath); os.IsNotExist(err) {
			// If the file doesn't exist, create a new one
			if err := os.MkdirAll(filepath.Dir(currentFilePath), 0666); err != nil {
				log.Fatalf("error creating directory: %v", err)
			}
			if _, err := os.Create(currentFilePath); err != nil {
				log.Fatalf("error creating file: %v", err)
			}
		}

		// Open the file with read/write, create or append mode
		var err error
		logFile, err = os.OpenFile(currentFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}

		logFilePath = currentFilePath
	}

	return logFile
}

// LogToFile creates a logger middleware that includes the API response message
func LogToFile() fiber.Handler {
	// Get the log file
	getLogFile()

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
		file := getLogFile()
		if _, err := file.WriteString(logEntry); err != nil {
			log.Printf("Error writing to log file: %v", err)
		}

		return err
	}
}

// LogError logs an error message to the log file
func LogError(c *fiber.Ctx, errorMsg string) {
	// Get the log file
	file := getLogFile()

	// Create a log entry for the error
	timestamp := time.Now().Format(time.RFC3339)
	statusCode := "500" // Internal Server Error
	ip := "-"
	method := "-"
	path := "-"

	if c != nil {
		ip = c.IP()
		method = c.Method()
		path = c.Path()
	}

	logEntry := timestamp +
		" | " + statusCode +
		" | ERROR | " + ip +
		" | " + method +
		" | " + path +
		" | " + errorMsg + "\n"

	// Append the log entry to the file
	logMutex.Lock()
	defer logMutex.Unlock()

	if _, err := file.WriteString(logEntry); err != nil {
		log.Printf("Error writing error to log file: %v", err)
	}
}
