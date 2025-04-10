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
func generateJWT(email string, userID int) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", errors.New("missing JWT secret")
	}

	claims := jwt.MapClaims{
		"email":   email,
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
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

	token, err := generateJWT(user.Email, int(user.ID))
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
