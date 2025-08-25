package helper

import (
	"byu-crm-service/modules/auth/repository"
	"byu-crm-service/modules/auth/service"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AdminHandler struct {
	db *gorm.DB
}

type AdminLoginRequest struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

func NewAdminHandler(db *gorm.DB) *AdminHandler {
	return &AdminHandler{db: db}
}

// ShowLogin displays the admin login page
func (h *AdminHandler) ShowLogin(c *fiber.Ctx) error {
	// Check if already authenticated
	if token := c.Cookies("admin_token"); token != "" {
		if h.validateAdminRole(token) {
			return c.Redirect("/admin/dashboard")
		}
	}

	return c.SendFile("./static/admin-login.html")
}

// HandleLogin processes admin login using existing auth service
func (h *AdminHandler) HandleLogin(c *fiber.Ctx) error {
	var req AdminLoginRequest

	// Parse form data or JSON
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request format",
		})
	}

	// Create auth service directly
	authRepo := h.createAuthRepository()
	authService := h.createAuthService(authRepo)

	// Use the auth service to login
	token, err := authService.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// Validate that user has admin role
	if !h.validateAdminRole(token["access_token"]) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Access denied. Admin privileges required.",
		})
	}

	// Set secure HTTP-only cookie with the JWT token
	c.Cookie(&fiber.Cookie{
		Name:     "admin_token",
		Value:    token["access_token"],
		Expires:  time.Now().Add(24 * time.Hour), // 24 hours
		HTTPOnly: true,
		Secure:   os.Getenv("APP_ENV") == "production",
		SameSite: "Lax",
	})

	// Return success response
	return c.JSON(fiber.Map{
		"status":   "success",
		"message":  "Login successful",
		"redirect": "/admin/dashboard",
	})
}

// HandleLogout processes admin logout
func (h *AdminHandler) HandleLogout(c *fiber.Ctx) error {
	// Clear the admin token cookie
	c.Cookie(&fiber.Cookie{
		Name:     "admin_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // Expire immediately
		HTTPOnly: true,
		Secure:   os.Getenv("APP_ENV") == "production",
		SameSite: "Lax",
	})

	return c.Redirect("/admin/login")
}

// ShowDashboard displays the admin dashboard
func (h *AdminHandler) ShowDashboard(c *fiber.Ctx) error {
	return c.SendFile("./static/admin-dashboard.html")
}

// createAuthRepository creates an auth repository instance
func (h *AdminHandler) createAuthRepository() repository.AuthRepository {
	return repository.NewAuthRepository(h.db)
}

// createAuthService creates an auth service instance
func (h *AdminHandler) createAuthService(authRepo repository.AuthRepository) service.AuthService {
	return service.NewAuthService(authRepo)
}

// validateAdminRole checks if the JWT token contains admin role
func (h *AdminHandler) validateAdminRole(tokenString string) bool {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return false
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check if user has Super-Admin role
		if userRole, ok := claims["user_role"].(string); ok {
			return userRole == "Super-Admin"
		}
	}

	return false
}

// AdminAuthMiddleware protects admin routes
func AdminAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from cookie
		token := c.Cookies("admin_token")
		if token == "" {
			return c.Redirect("/admin/login")
		}

		// Validate token and check admin role
		adminHandler := &AdminHandler{}
		if !adminHandler.validateAdminRole(token) {
			// Clear invalid cookie
			c.Cookie(&fiber.Cookie{
				Name:     "admin_token",
				Value:    "",
				Expires:  time.Now().Add(-time.Hour),
				HTTPOnly: true,
				Secure:   os.Getenv("APP_ENV") == "production",
				SameSite: "Lax",
			})
			return c.Redirect("/admin/login")
		}

		// Set admin context and user info from JWT
		c.Locals("admin", true)

		// Parse JWT to extract user info for context
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret != "" {
			if parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			}); err == nil {
				if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
					if userID, ok := claims["user_id"].(float64); ok {
						c.Locals("user_id", int(userID))
					}
					if userRole, ok := claims["user_role"].(string); ok {
						c.Locals("user_role", userRole)
					}
					if email, ok := claims["email"].(string); ok {
						c.Locals("user_email", email)
					}
				}
			}
		}

		return c.Next()
	}
}
