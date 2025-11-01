package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"byu-crm-service/models"
	"byu-crm-service/modules/auth/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type authService struct {
	userRepo repository.AuthRepository
}

func NewAuthService(userRepo repository.AuthRepository) AuthService {
	return &authService{userRepo: userRepo}
}

// Generate JWT Token
func generateJWT(email string, userID int, userRole string, territoryType string, territoryID int, user_permissions []string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", errors.New("missing JWT secret")
	}

	claims := jwt.MapClaims{
		"email":          email,
		"user_id":        userID,
		"user_role":      userRole,
		"territory_type": territoryType,
		"territory_id":   territoryID,
		"permissions":    user_permissions,
		"exp":            time.Now().Add(time.Hour * 24 * 30).Unix(), // TTL: 1 bulan
		"iat":            time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (s *authService) GenerateAccessToken(email string, userID int, userRole string, territoryType string, territoryID int, user_permissions []string, adminID int) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", errors.New("missing JWT secret")
	}

	claims := jwt.MapClaims{
		"email":          email,
		"user_id":        userID,
		"user_role":      userRole,
		"territory_type": territoryType,
		"territory_id":   territoryID,
		"permissions":    user_permissions,
		"admin_id":       adminID,
		"exp":            time.Now().Add(30 * time.Minute).Unix(), // 🔥 access token 2 meni	t
		"iat":            time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func generateRefreshToken(adminID, userID int) (string, error) {
	refreshSecret := os.Getenv("REFRESH_SECRET") // 🔥 tambahkan secret baru
	if refreshSecret == "" {
		return "", errors.New("missing REFRESH secret")
	}

	claims := jwt.MapClaims{
		"user_id":  userID,
		"admin_id": adminID,
		"exp":      time.Now().Add(time.Hour * 24 * 30).Unix(), // 🔥 refresh token 30 hari
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(refreshSecret))
}

// 🔥 ubah return: sekarang return accessToken + refreshToken
func (s *authService) Login(email, password string) (map[string]string, error) {
	user, err := s.userRepo.GetUserByKey("email", email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !s.userRepo.CheckPassword(password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	// role check tetap
	allowed := true
	for _, role := range user.RoleNames {
		if role == "Super-Admin" || role == "YAE" {
			allowed = true
			break
		}
	}
	if !allowed {
		return nil, errors.New("you cannot access this application")
	}

	// 🔥 generate access + refresh token
	accessToken, err := s.GenerateAccessToken(user.Email, int(user.ID), user.RoleNames[0], user.TerritoryType, int(user.TerritoryID), user.Permissions, 0)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	refreshToken, err := generateRefreshToken(0, int(user.ID))
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	return map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil
}

func (s *authService) Impersonate(adminID int, email string) (map[string]string, error) {
	user, err := s.userRepo.GetUserByKey("email", email)
	if err != nil {
		return nil, errors.New("invalid email")
	}

	// 🔥 generate access + refresh token
	accessToken, err := s.GenerateAccessToken(user.Email, int(user.ID), user.RoleNames[0], user.TerritoryType, int(user.TerritoryID), user.Permissions, adminID)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	refreshToken, err := generateRefreshToken(adminID, int(user.ID))
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	return map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil
}

func (s *authService) CheckPassword(password, hashedPassword string) bool {
	return s.userRepo.CheckPassword(password, hashedPassword)
}

func (s *authService) GetUserByKey(key, value string) (*models.User, error) {
	return s.userRepo.GetUserByKey(key, value)
}

// getGoogleOAuthConfig returns the OAuth2 config for Google
func getGoogleOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URI"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

// GetGoogleOAuthURL returns the URL to redirect the user to for Google OAuth
func (s *authService) GetGoogleOAuthURL() string {
	config := getGoogleOAuthConfig()
	return config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
}

// HandleGoogleCallback handles the callback from Google OAuth
func (s *authService) HandleGoogleCallback(code string) (string, error) {
	config := getGoogleOAuthConfig()

	// Exchange the authorization code for a token
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		return "", fmt.Errorf("code exchange failed: %s", err.Error())
	}

	// Get user info from Google
	userInfo, err := getUserInfoFromGoogle(token.AccessToken)
	if err != nil {
		return "", fmt.Errorf("failed to get user info: %s", err.Error())
	}

	// Check if user exists in our database
	user, err := s.userRepo.GetUserByKey("email", userInfo["email"].(string))
	if err != nil {
		// User doesn't exist, create a new user
		// For now, we'll just return an error
		return "", errors.New("user not found in our system")
	}

	// Update Google ID if not already set
	if user.GoogleID == "" {
		userMap := map[string]interface{}{
			"google_id": userInfo["id"].(string),
		}
		if _, err := s.userRepo.UpdateUser(int(user.ID), userMap); err != nil {
			return "", fmt.Errorf("failed to update user: %s", err.Error())
		}
	}

	// Check if user has the required roles
	allowed := false
	for _, role := range user.RoleNames {
		if role == "Super-Admin" || role == "YAE" {
			allowed = true
			break
		}
	}

	if !allowed {
		return "", errors.New("you cannot access this application")
	}

	// Generate JWT token
	jwtToken, err := generateJWT(user.Email, int(user.ID), user.RoleNames[0], user.TerritoryType, int(user.TerritoryID), user.Permissions)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %s", err.Error())
	}

	return jwtToken, nil
}

// getUserInfoFromGoogle gets the user info from Google using the access token
func getUserInfoFromGoogle(accessToken string) (map[string]interface{}, error) {
	// Make a request to the Google API to get user info
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse the JSON response
	var userInfo map[string]interface{}
	if err := json.Unmarshal(data, &userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}
