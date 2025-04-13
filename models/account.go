package models

import "time"

type Account struct {
	ID                      uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	AccountImage            *string   `json:"account_image"`
	AccountName             *string   `json:"account_name"`
	AccountType             *string   `json:"account_type"`
	AccountCategory         *string   `json:"account_category"`
	AccountCode             *string   `json:"account_code"`
	City                    *uint     `json:"city"`
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

	SocialMedias               []SocialMedia               `json:"social_medias" gorm:"foreignKey:SubjectID;references:ID"`
	AccountTypeCampusDetail    *AccountTypeCampusDetail    `json:"campus_detail,omitempty" gorm:"foreignKey:AccountID"`
	AccountTypeSchoolDetail    *AccountTypeSchoolDetail    `json:"school_detail,omitempty" gorm:"foreignKey:AccountID"`
	AccountTypeCommunityDetail *AccountTypeCommunityDetail `json:"community_detail,omitempty" gorm:"foreignKey:AccountID"`
	AccountCity                *City                       `json:"account_city" gorm:"foreignKey:City;references:ID"`
	AccountFaculties           []AccountFaculty            `gorm:"foreignKey:AccountID" json:"account_faculties"`
}
