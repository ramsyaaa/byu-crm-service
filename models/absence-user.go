package models

import "time"

type AbsenceUser struct {
	ID                       uint                      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID                   *uint                     `json:"user_id"`
	SubjectType              *string                   `json:"subject_type"`
	SubjectID                *uint                     `json:"subject_id"`
	Type                     *string                   `json:"type"`
	ClockIn                  time.Time                 `json:"clock_in"`
	ClockOut                 *time.Time                `json:"clock_out"`
	Description              string                    `json:"description"`
	Longitude                string                    `json:"longitude"`
	Latitude                 string                    `json:"latitude"`
	Status                   *uint                     `json:"status"`
	EvidenceImage            *string                   `json:"evidence_image"`
	CreatedAt                time.Time                 `json:"created_at"`
	UpdatedAt                time.Time                 `json:"updated_at"`
	Account                  *Account                  `gorm:"-" json:"account"`
	VisitHistory             *VisitHistory             `gorm:"-" json:"visit_history"`
	UserName                 string                    `gorm:"->" json:"user_name"`
	YaeCode                  *string                   `gorm:"->" json:"yae_code"`
	AccountTypeSchoolDetails *AccountTypeSchoolDetails `json:"account_type_school_details,omitempty" gorm:"-"`
}
