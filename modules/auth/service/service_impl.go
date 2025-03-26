package service

import (
	"errors"
	"os"
	"time"

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
func generateJWT(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := os.Getenv("JWT_SECRET")
	return token.SignedString([]byte(jwtSecret))
}

// Login Service
func (s *authService) Login(email, password string) (string, error) {
	hashedPassword, err := s.userRepo.GetUserByKey("email", email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	dbPassword := hashedPassword.Password

	if !s.userRepo.CheckPassword(password, dbPassword) {
		return "", errors.New("invalid email or password")
	}

	token, err := generateJWT(email)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}
