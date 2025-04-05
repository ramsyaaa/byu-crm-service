package models

import "time"

type AccountMember struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	SubjectType *string   `json:"subject_type"`
	SubjectID   *uint     `json:"subject_id"`
	Year        *string   `json:"year"`
	Amount      *string   `json:"amount"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
