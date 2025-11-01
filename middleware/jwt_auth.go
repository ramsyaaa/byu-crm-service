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

		adminID, ok := claims["admin_id"].(float64)
		if !ok {
			adminID = 0
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
		c.Locals("admin_id", int(adminID))

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

func unauthorized(c *fiber.Ctx, message string) error {
	response := helper.APIResponse(message, fiber.StatusUnauthorized, "error", nil)
	return c.Status(fiber.StatusUnauthorized).JSON(response)
}
