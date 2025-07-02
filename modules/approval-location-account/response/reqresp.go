package response

import (
	"byu-crm-service/models"
	"time"
)

type ApprovalLocationAccountResponse struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    *uint     `json:"user_id"`
	AccountID *uint     `json:"account_id"`
	Longitude *string   `json:"longitude"`
	Latitude  *string   `json:"latitude"`
	Status    *uint     `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	LatestLongitude *string `json:"latest_longitude"` // Longitude terakhir yang diset
	LatestLatitude  *string `json:"latest_latitude"`  // Latitude terakhir yang dis

	AccountName *string `json:"account_name"` // Nama akun untuk ditampilkan
	UserName    *string `json:"user_name"`    // Nama pengguna untuk ditampilkan

	User    *models.User    `json:"user"`   // Relasi user
	Account *models.Account `jaon:"acount"` // Relasi account
}
