package middleware

import (
	"byu-crm-service/helper"
	"fmt"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
)

// RecoverMiddleware catches any panics and logs them to the logging system
func RecoverMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				// Get stack trace
				stackTrace := debug.Stack()

				// Log the panic with stack trace
				errorMsg := fmt.Sprintf("PANIC RECOVERED: %v\nStack Trace:\n%s", r, stackTrace)

				// Log to file using our enhanced logger's format
				helper.LogError(c, errorMsg)

				// Return a 500 response to the client
				// Check if headers have been sent by checking if status code is not 0
				if c.Response().StatusCode() == 0 {
					c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"status":  "error",
						"message": "Internal Server Error",
					})
				}
			}
		}()

		return c.Next()
	}
}
