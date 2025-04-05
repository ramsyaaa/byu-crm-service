package middleware

import (
	"byu-crm-service/helper"
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(c *fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
		ContextKey: "jwt", // default: "user"
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			response := helper.APIResponse(err.Error(), fiber.StatusUnauthorized, "error", nil)
			return c.Status(fiber.StatusUnauthorized).JSON(response)
		},
	})(c)
}

// JWTUserContextMiddleware extract user_id and store in c.Locals
func JWTUserContextMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("jwt") // from jwtware.Config.ContextKey
		claims, ok := user.(*jwt.Token)
		if !ok {
			response := helper.APIResponse("Unauthorized: Invalid token format", fiber.StatusUnauthorized, "error", nil)
			return c.Status(fiber.StatusUnauthorized).JSON(response)
		}

		mapClaims, ok := claims.Claims.(jwt.MapClaims)
		if !ok {
			response := helper.APIResponse("Unauthorized: Invalid claims", fiber.StatusUnauthorized, "error", nil)
			return c.Status(fiber.StatusUnauthorized).JSON(response)
		}

		userIDFloat, ok := mapClaims["user_id"].(float64)
		if !ok {
			response := helper.APIResponse("Unauthorized: Invalid user ID", fiber.StatusUnauthorized, "error", nil)
			return c.Status(fiber.StatusUnauthorized).JSON(response)
		}

		// Simpan ke context
		c.Locals("user_id", int(userIDFloat))

		return c.Next()
	}
}
