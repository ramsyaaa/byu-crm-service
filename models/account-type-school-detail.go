package models

import "time"

type AccountTypeSchoolDetail struct {
	ID                      uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	AccountID               *uint     `json:"account_id"`
	DiesNatalis             time.Time `json:"dies_natalis"`
	Extracurricular         *string   `json:"extracurricular"`
	FootballFieldBrannnding *string   `json:"football_field_brannnding"`
	BasketballFieldBranding *string   `json:"basketball_field_branding"`
	WallPaintingBranding    *string   `json:"wall_painting_branding"`
	WallMagazineBranding    *string   `json:"wall_magazine_branding"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}
