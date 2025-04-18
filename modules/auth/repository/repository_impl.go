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

	// Ambil user berdasarkan key
	err := r.db.Where(key+" = ?", value).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	}

	// Ambil role_id dari tabel model_has_roles
	var roleIDs []uint
	if err := r.db.Table("model_has_roles").
		Where("model_id = ? AND model_type = ?", user.ID, "App\\Models\\User").
		Pluck("role_id", &roleIDs).Error; err != nil {
		return nil, err
	}

	// Ambil nama role dari tabel roles
	if len(roleIDs) > 0 {
		var roleNames []string
		if err := r.db.Table("roles").
			Where("id IN ?", roleIDs).
			Pluck("name", &roleNames).Error; err != nil {
			return nil, err
		}
		user.RoleNames = roleNames
	}

	return &user, nil
}

// Cek apakah password cocok
func (r *authRepository) CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
