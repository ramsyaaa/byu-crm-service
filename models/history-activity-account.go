package models

import "time"

type HistoryActivityAccount struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      uint      `gorm:"column:user_id"`
	AccountID   uint      `gorm:"column:account_id"`
	Type        string    `gorm:"column:type"`
	SubjectType *string   `gorm:"column:subject_type"`
	SubjectID   *uint     `gorm:"column:subject_id"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}
