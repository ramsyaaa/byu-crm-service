package models

import "time"

type SocialMedia struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	SubjectType *string   `json:"subject_type"`
	SubjectID   *uint     `json:"subject_id"`
	Category    *string   `json:"category"`
	Url         *string   `json:"url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
