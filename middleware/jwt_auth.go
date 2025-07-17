package middleware

import (
	"byu-crm-service/helper"
	"os"
	"strings"

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

// AdminJWTMiddleware handles JWT authentication specifically for admin routes
func AdminJWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return handleAdminAuthError(c, "Unauthorized: missing authorization header")
		}

		// Check if it starts with "Bearer "
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			return handleAdminAuthError(c, "Unauthorized: invalid authorization format")
		}

		// Extract token
		tokenString := authHeader[7:]
		if tokenString == "" {
			return handleAdminAuthError(c, "Unauthorized: missing JWT token")
		}

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid signing method")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			return handleAdminAuthError(c, "Unauthorized: invalid JWT token - "+err.Error())
		}

		if !token.Valid {
			return handleAdminAuthError(c, "Unauthorized: invalid JWT token")
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return handleAdminAuthError(c, "Unauthorized: invalid token claims")
		}

		// Extract user email
		email, ok := claims["email"].(string)
		if !ok {
			return handleAdminAuthError(c, "Unauthorized: email not found in token")
		}

		// Check user_role from JWT token
		userRole, ok := claims["user_role"].(string)
		if !ok {
			return handleAdminAuthError(c, "Unauthorized: user role not found in token")
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
		c.Locals("jwt", token) // Store token for compatibility

		return c.Next()
	}
}

// handleAdminAuthError handles authentication errors for admin routes
func handleAdminAuthError(c *fiber.Ctx, message string) error {
	path := c.Path()

	// For API requests or admin profile endpoint, return JSON error
	if c.Get("Accept") == "application/json" ||
		path == "/admin/profile" ||
		path == "/api-logs" ||
		strings.HasPrefix(path, "/api-logs/") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"meta": fiber.Map{
				"status":  "error",
				"message": message,
				"code":    fiber.StatusUnauthorized,
			},
			"data": nil,
		})
	}

	// For browser requests to dashboard, redirect to login (but avoid infinite loops)
	if path != "/admin/login" {
		return c.Redirect("/admin/login")
	}

	// If we're already on login page, return JSON error to avoid loops
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"meta": fiber.Map{
			"status":  "error",
			"message": message,
			"code":    fiber.StatusUnauthorized,
		},
		"data": nil,
	})
}

// AdminAuthMiddleware checks if user is authenticated and has Super-Admin user_role
// This is now a secondary middleware that assumes JWT is already validated by AdminJWTMiddleware
func AdminAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if user info is already set by AdminJWTMiddleware
		email := c.Locals("admin_user_email")
		userRole := c.Locals("admin_user_role")

		if email == nil || userRole == nil {
			return handleAdminAuthError(c, "Unauthorized: admin authentication required")
		}

		// Additional role validation (redundant but safe)
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

		return c.Next()
	}
}

func unauthorized(c *fiber.Ctx, message string) error {
	response := helper.APIResponse(message, fiber.StatusUnauthorized, "error", nil)
	return c.Status(fiber.StatusUnauthorized).JSON(response)
}
