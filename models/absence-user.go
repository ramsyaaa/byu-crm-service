package models

import "time"

type AbsenceUser struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      *uint     `json:"user_id"`
	SubjectType *string   `json:"subject_type"`
	SubjectID   *uint     `json:"subject_id"`
	Type        *string   `json:"type"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Longitude   string    `json:"longitude"`
	Latitude    string    `json:"latitude"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
