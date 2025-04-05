package models

import "time"

type AccountTypeCampusDetail struct {
	ID                       uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	AccountID                *uint     `json:"account_id"`
	RangeAge                 *string   `json:"range_age"`
	Origin                   *string   `gorm:"type:longtext" json:"origin"`
	OrganizationName         *string   `gorm:"type:longtext" json:"organization_name"`
	PreferenceTechnologies   *string   `gorm:"type:longtext" json:"preference_technologies"`
	MemberNeeds              *string   `gorm:"type:longtext" json:"member_needs"`
	ItInfrastructures        *string   `gorm:"type:longtext" json:"it_infrastructures"`
	DigitalCollaborations    *string   `gorm:"type:longtext" json:"digital_collaborations"`
	Byod                     *uint     `json:"byod"`
	AccessTechnology         *string   `gorm:"type:longtext" json:"access_technology"`
	CampusAdministrationApp  *string   `gorm:"type:longtext" json:"campus_administration_app"`
	PotentionalCollaboration *string   `gorm:"type:longtext" json:"potentional_collaboration"`
	UniversityRank           *string   `gorm:"type:longtext" json:"university_rank"`
	FocusProgramStudy        *string   `gorm:"type:longtext" json:"focus_program_study"`
	ProgramIdentification    *string   `gorm:"type:longtext" json:"program_identification"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`
}
