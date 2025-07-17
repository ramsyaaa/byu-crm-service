package middleware

import (
	"byu-crm-service/helper"
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(c *fiber.Ctx) error {
	// Jalankan middleware jwtware
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
		ContextKey:   "jwt",
		TokenLookup:  "header:Authorization", // tetap sama
		AuthScheme:   "Bearer",
		ErrorHandler: jwtErrorHandler,
	})(c)
}

// JWTMiddlewareHandler returns a JWT middleware handler
func JWTMiddlewareHandler() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
		ContextKey:   "jwt",
		TokenLookup:  "header:Authorization",
		AuthScheme:   "Bearer",
		ErrorHandler: jwtErrorHandler,
	})
}

func jwtErrorHandler(c *fiber.Ctx, err error) error {
	response := helper.APIResponse("Unauthorized: "+err.Error(), fiber.StatusUnauthorized, "error", nil)
	return c.Status(fiber.StatusUnauthorized).JSON(response)
}

// JWTUserContextMiddleware extract user_id and store in c.Locals
func JWTUserContextMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("jwt")
		token, ok := user.(*jwt.Token)
		if !ok || token == nil {
			return unauthorized(c, "Unauthorized: Invalid token format")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return unauthorized(c, "Unauthorized: Invalid claims")
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			return unauthorized(c, "Unauthorized: user id not found in token")
		}

		userRole, ok := claims["user_role"].(string)
		if !ok {
			return unauthorized(c, "Unauthorized: user role not found in token")
		}

		territoryType, ok := claims["territory_type"].(string)
		if !ok {
			return unauthorized(c, "Unauthorized: territory type not found in token")
		}

		territoryID, ok := claims["territory_id"].(float64)
		if !ok {
			return unauthorized(c, "Unauthorized: territory id not found in token")
		}

		// Ambil permission
		permissions, ok := claims["permissions"].([]interface{})
		if !ok {
			return unauthorized(c, "Unauthorized: permissions not found or invalid")
		}

		var permStrings []string
		for _, p := range permissions {
			if str, ok := p.(string); ok {
				permStrings = append(permStrings, str)
			}
		}

		// Set ke context
		c.Locals("user_id", int(userID))
		c.Locals("user_role", userRole)
		c.Locals("territory_type", territoryType)
		c.Locals("territory_id", int(territoryID))
		c.Locals("permissions", permStrings)

		return c.Next()
	}
}

func PermissionMiddleware(requiredPermission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		perms, ok := c.Locals("permissions").([]string)
		if !ok {
			return jwtErrorHandler(c, fiber.NewError(fiber.StatusUnauthorized, "Permission context invalid"))
		}

		for _, p := range perms {
			if p == requiredPermission {
				return c.Next()
			}
		}
		return jwtErrorHandler(c, fiber.NewError(fiber.StatusUnauthorized, "You cannot access "+requiredPermission))
	}
}

// AdminAuthMiddleware checks if user is authenticated and has Super-Admin user_role
func AdminAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// First check if user is authenticated via JWT
		user := c.Locals("jwt")
		token, ok := user.(*jwt.Token)
		if !ok || token == nil {
			// For API requests, return JSON error instead of redirect
			if c.Get("Accept") == "application/json" || c.Path() != "/admin/dashboard" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"meta": fiber.Map{
						"status":  "error",
						"message": "Unauthorized: missing or malformed JWT",
						"code":    fiber.StatusUnauthorized,
					},
					"data": nil,
				})
			}
			// Redirect to login page for admin interface
			return c.Redirect("/admin/login")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			if c.Get("Accept") == "application/json" || c.Path() != "/admin/dashboard" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"meta": fiber.Map{
						"status":  "error",
						"message": "Unauthorized: invalid token claims",
						"code":    fiber.StatusUnauthorized,
					},
					"data": nil,
				})
			}
			return c.Redirect("/admin/login")
		}

		// Extract user email
		email, ok := claims["email"].(string)
		if !ok {
			if c.Get("Accept") == "application/json" || c.Path() != "/admin/dashboard" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"meta": fiber.Map{
						"status":  "error",
						"message": "Unauthorized: email not found in token",
						"code":    fiber.StatusUnauthorized,
					},
					"data": nil,
				})
			}
			return c.Redirect("/admin/login")
		}

		// Check user_role from JWT token
		userRole, ok := claims["user_role"].(string)
		if !ok {
			if c.Get("Accept") == "application/json" || c.Path() != "/admin/dashboard" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"meta": fiber.Map{
						"status":  "error",
						"message": "Unauthorized: user role not found in token",
						"code":    fiber.StatusUnauthorized,
					},
					"data": nil,
				})
			}
			return c.Redirect("/admin/login")
		}

		// Check if user has Super-Admin privileges
		if userRole != "Super-Admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"meta": fiber.Map{
					"status":  "error",
					"message": "Access denied. Super-Admin privileges required.",
					"code":    fiber.StatusForbidden,
				},
				"data": nil,
			})
		}

		// Store user info in context for use in handlers
		c.Locals("admin_user_email", email)
		c.Locals("admin_user_role", userRole)

		return c.Next()
	}
}

func unauthorized(c *fiber.Ctx, message string) error {
	response := helper.APIResponse(message, fiber.StatusUnauthorized, "error", nil)
	return c.Status(fiber.StatusUnauthorized).JSON(response)
}
