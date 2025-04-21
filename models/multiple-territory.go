package models

type MultipleTerritory struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	SubjectType string `json:"subject_type"`
	SubjectIDs  string `json:"subject_ids"`
	UserID      int    `json:"user_id"`
}
