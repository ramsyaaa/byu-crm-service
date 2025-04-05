package repository

import (
	"errors"

	"byu-crm-service/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) GetUserByKey(key, value string) (*models.User, error) {
	var user models.User
	err := r.db.Where(key+" = ?", value).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

// Cek apakah password cocok
func (r *authRepository) CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
