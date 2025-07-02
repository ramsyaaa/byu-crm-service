package response

import "time"

type NotificationResponse struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	CallbackUrl *string   `json:"callback_url"`
	SubjectType *string   `json:"subject_type"`
	SubjectID   *uint     `json:"subject_id"`
	UserID      *uint     `json:"user_id"`
	IsRead      *uint     `json:"is_read"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
