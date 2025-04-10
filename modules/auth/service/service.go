package service

import "byu-crm-service/models"

type AuthService interface {
	Login(email, password string) (string, error)
	CheckPassword(password, hashedPassword string) bool
	GetUserByKey(key, value string) (*models.User, error)
}
