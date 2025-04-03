package models

import "time"

type AccountTypeCommunityDetail struct {
	ID                          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	AccountID                   *uint     `json:"account_id"`
	AccountSubtype              *string   `json:"account_subtype"`
	Group                       *string   `json:"group"`
	GroupName                   *string   `json:"group_name"`
	RangeAge                    *string   `gorm:"type:longtext" json:"range_age"`
	Gender                      *string   `gorm:"type:longtext" json:"gender"`
	EducationalBackground       *string   `gorm:"type:longtext" json:"educational_background"`
	Profession                  *string   `gorm:"type:longtext" json:"profession"`
	Income                      *string   `gorm:"type:longtext" json:"income"`
	ProductService              *string   `gorm:"type:longtext" json:"product_service"`
	PotentialCollaborationItems *string   `gorm:"type:longtext" json:"potential_collaboration_items"`
	PotentionalCollaboration    *string   `gorm:"type:longtext" json:"potentional_collaboration"`
	CreatedAt                   time.Time `json:"created_at"`
	UpdatedAt                   time.Time `json:"updated_at"`
}
