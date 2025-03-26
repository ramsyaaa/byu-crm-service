package repository

import "byu-crm-service/models"

type AuthRepository interface {
	GetUserByKey(key, value string) (*models.User, error)
	CheckPassword(password, hashedPassword string) bool
}
