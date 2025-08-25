package service

import "byu-crm-service/models"

type AuthService interface {
	Login(email, password string) (map[string]string, error)
	CheckPassword(password, hashedPassword string) bool
	GetUserByKey(key, value string) (*models.User, error)
	GenerateAccessToken(email string, userID int, userRole string, territoryType string, territoryID int, user_permissions []string) (string, error)
	// Google OAuth methods
	GetGoogleOAuthURL() string
	HandleGoogleCallback(code string) (string, error)
}
