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
	authHeader := c.Get("Authorization")

	if authHeader == "" {
		return jwtErrorHandler(c, fiber.NewError(fiber.StatusUnauthorized, "Authorization header tidak ditemukan"))
	}

	// Pisahkan berdasarkan spasi → [Bearer, Bearer, <token>]
	parts := strings.Fields(authHeader)
	if len(parts) < 2 {
		return jwtErrorHandler(c, fiber.NewError(fiber.StatusUnauthorized, "Format Authorization salah"))
	}

	// Ambil token paling akhir (menghindari Bearer double)
	tokenOnly := parts[len(parts)-1]
	c.Request().Header.Set("Authorization", tokenOnly) // Overwrite agar jwtware bisa proses

	// Lanjut ke JWT middleware
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
		ContextKey:   "jwt",
		TokenLookup:  "header:Authorization", // tetap gunakan header
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
			return unauthorized(c, "Unauthorized: user_id not found in token")
		}

		c.Locals("user_id", int(userID))
		return c.Next()
	}
}

func unauthorized(c *fiber.Ctx, message string) error {
	response := helper.APIResponse(message, fiber.StatusUnauthorized, "error", nil)
	return c.Status(fiber.StatusUnauthorized).JSON(response)
}
