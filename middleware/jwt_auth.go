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

		c.Locals("user_id", int(userID))
		c.Locals("user_role", userRole)
		c.Locals("territory_type", territoryType)
		c.Locals("territory_id", int(territoryID))
		return c.Next()
	}
}

func unauthorized(c *fiber.Ctx, message string) error {
	response := helper.APIResponse(message, fiber.StatusUnauthorized, "error", nil)
	return c.Status(fiber.StatusUnauthorized).JSON(response)
}
