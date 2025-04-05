package models

import "time"

type VisitHistory struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        *uint     `json:"user_id"`
	SubjectType   *string   `json:"subject_type"`
	SubjectID     *uint     `json:"subject_id"`
	AbsenceUserID *uint     `json:"absence_user_id"`
	Greeting      bool      `json:"greeting"`
	Survey        bool      `json:"survey"`
	Presentation  bool      `json:"presentation"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
