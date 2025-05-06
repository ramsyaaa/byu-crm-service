package response

import (
	"byu-crm-service/models"
	"time"
)

type TerritoryCategory struct {
	Category string `json:"category"`
	Count    int64  `json:"count"`
}

// Struktur untuk hasil wilayah dengan kategori
type TerritoryResult struct {
	ID         int64               `json:"id"`
	Name       string              `json:"name"`
	Total      int64               `json:"total"`
	Categories []TerritoryCategory `json:"categories"`
}

type AccountResponse struct {
	ID                      uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	AccountName             *string   `json:"account_name"`
	AccountType             *string   `json:"account_type"`
	AccountCategory         *string   `json:"account_category"`
	AccountCode             *string   `json:"account_code"`
	City                    *uint     `json:"city"`
	CityName                *string   `json:"city_name" gorm:"column:city_name"`
	ClusterName             *string   `json:"cluster_name" gorm:"column:cluster_name"`
	BranchName              *string   `json:"branch_name" gorm:"column:branch_name"`
	RegionName              *string   `json:"region_name" gorm:"column:region_name"`
	AreaName                *string   `json:"area_name" gorm:"column:area_name"`
	ContactName             *string   `json:"contact_name"`
	EmailAccount            *string   `json:"email_account"`
	WebsiteAccount          *string   `json:"website_account"`
	Potensi                 *string   `json:"potensi"`
	SystemInformasiAkademik *string   `json:"system_informasi_akademik"`
	CustomerSegmentationId  *string   `json:"customer_segmentation_id"`
	Latitude                *string   `json:"latitude"`
	Longitude               *string   `json:"longitude"`
	Ownership               *string   `json:"ownership"`
	Pic                     *string   `json:"pic"`
	PicInternal             *string   `json:"pic_internal"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`

	SocialMedias               []models.SocialMedia               `json:"social_medias" gorm:"foreignKey:SubjectID;references:ID"`
	AccountTypeCampusDetail    *models.AccountTypeCampusDetail    `json:"campus_detail,omitempty" gorm:"foreignKey:AccountID"`
	AccountTypeSchoolDetail    *models.AccountTypeSchoolDetail    `json:"school_detail,omitempty" gorm:"foreignKey:AccountID"`
	AccountTypeCommunityDetail *models.AccountTypeCommunityDetail `json:"community_detail,omitempty" gorm:"foreignKey:AccountID"`
	AccountCity                *models.City                       `json:"account_city" gorm:"foreignKey:City;references:ID"`
	AccountFaculties           []models.AccountFaculty            `gorm:"foreignKey:AccountID" json:"account_faculties"`
	AccountMembers             []models.AccountMember             `json:"account_members" gorm:"foreignKey:SubjectID;references:ID;where:subject_type='App\\Models\\Account'"`
	AccountLectures            []models.AccountMember             `json:"account_lectures" gorm:"foreignKey:SubjectID;references:ID;where:subject_type='App\\Models\\AccountLecture'"`
	Contacts                   []models.Contact                   `gorm:"many2many:contact_accounts;foreignKey:ID;joinForeignKey:account_id;References:ID;joinReferences:contact_id" json:"contacts"`
}

type SingleAccountResponse struct {
	ID                      uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	AccountName             *string   `json:"account_name"`
	AccountType             *string   `json:"account_type"`
	AccountCategory         *string   `json:"account_category"`
	AccountCode             *string   `json:"account_code"`
	City                    *uint     `json:"city"`
	CityName                *string   `json:"city_name" gorm:"column:city_name"`
	ClusterName             *string   `json:"cluster_name" gorm:"column:cluster_name"`
	BranchName              *string   `json:"branch_name" gorm:"column:branch_name"`
	RegionName              *string   `json:"region_name" gorm:"column:region_name"`
	AreaName                *string   `json:"area_name" gorm:"column:area_name"`
	ContactName             *string   `json:"contact_name"`
	EmailAccount            *string   `json:"email_account"`
	WebsiteAccount          *string   `json:"website_account"`
	Potensi                 *string   `json:"potensi"`
	SystemInformasiAkademik *string   `json:"system_informasi_akademik"`
	CustomerSegmentationId  *string   `json:"customer_segmentation_id"`
	Latitude                *string   `json:"latitude"`
	Longitude               *string   `json:"longitude"`
	Ownership               *string   `json:"ownership"`
	Pic                     *string   `json:"pic"`
	PicInternal             *string   `json:"pic_internal"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`

	Contacts  []models.Contact `json:"contacts"`
	ContactID []string         `json:"contact_id"`

	Url      []string `json:"url"`
	Category []string `json:"category"`

	// only for account category school
	DiesNatalis             *time.Time `json:"dies_natalis"`
	Extracurricular         *string    `json:"extracurricular"`
	FootballFieldBranding   *string    `json:"football_field_branding"`
	BasketballFieldBranding *string    `json:"basketball_field_branding"`
	WallPaintingBranding    *string    `json:"wall_painting_branding"`
	WallMagazineBranding    *string    `json:"wall_magazine_branding"`

	// account category campus
	Faculties               []string `json:"faculties"`
	YearLecture             []string `json:"year_lecture"`
	AmountLecture           []string `json:"amount_lecture"`
	Origin                  []string `json:"origin"`
	PercentageOrigin        []string `json:"percentage_origin"`
	OrganizationName        []string `json:"organization_name"`
	PreferenceTechnologies  []string `json:"preference_technologies"`
	MemberNeeds             []string `json:"member_needs"`
	AccessTechnology        *string  `json:"access_technology"`
	Byod                    *string  `json:"byod"`
	ItInfrastructures       []string `json:"it_infrastructures"`
	DigitalCollaborations   []string `json:"digital_collaborations"`
	CampusAdministrationApp *string  `json:"campus_administration_app"`
	ProgramIdentification   []string `json:"program_identification"`
	YearRank                []string `json:"year_rank"`
	Rank                    []string `json:"rank"`
	ProgramStudy            []string `json:"program_study"`

	// account category campus & community
	Year                     []string `json:"year"`
	Amount                   []string `json:"amount"`
	Age                      []string `json:"age"`
	PercentageAge            []string `json:"percentage_age"`
	ScheduleCategory         []string `json:"schedule_category"`
	Title                    []string `json:"title"`
	Date                     []string `json:"date"`
	PotentionalCollaboration *string  `json:"potentional_collaboration"`

	// account category community
	AccountSubtype                  *string  `json:"account_subtype"`
	Group                           *string  `json:"group"`
	GroupName                       *string  `json:"group_name"`
	ProductService                  *string  `json:"product_service"`
	PotentialCollaborationItems     *string  `json:"potential_collaboration_items"`
	Gender                          []string `json:"gender"`
	PercentageGender                []string `json:"percentage_gender"`
	EducationalBackground           []string `json:"educational_background"`
	PercentageEducationalBackground []string `json:"percentage_educational_background"`
	Profession                      []string `json:"profession"`
	PercentageProfession            []string `json:"percentage_profession"`
	Income                          []string `json:"income"`
	PercentageIncome                []string `json:"percentage_income"`

	SocialMedias []models.SocialMedia `json:"social_medias" gorm:"foreignKey:SubjectID;references:ID"`
}
