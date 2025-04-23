package service

import (
	"errors"
	"os"
	"time"

	"byu-crm-service/models"
	"byu-crm-service/modules/auth/repository"

	"github.com/golang-jwt/jwt/v5"
)

type authService struct {
	userRepo repository.AuthRepository
}

func NewAuthService(userRepo repository.AuthRepository) AuthService {
	return &authService{userRepo: userRepo}
}

// Generate JWT Token
func generateJWT(email string, userID int, userRole string, territoryType string, territoryID int) (string, error) {
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

// Login Service
func (s *authService) Login(email, password string) (string, error) {
	user, err := s.userRepo.GetUserByKey("email", email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	if !s.userRepo.CheckPassword(password, user.Password) {
		return "", errors.New("invalid email or password")
	}

	//  Check if user has the required roles
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

	token, err := generateJWT(user.Email, int(user.ID), user.RoleNames[0], user.TerritoryType, int(user.TerritoryID))
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}

func (s *authService) CheckPassword(password, hashedPassword string) bool {
	return s.userRepo.CheckPassword(password, hashedPassword)
}

func (s *authService) GetUserByKey(key, value string) (*models.User, error) {
	return s.userRepo.GetUserByKey(key, value)
}
