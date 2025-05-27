package models

import "time"

type Account struct {
	ID                      uint      `gorm:"primaryKey;autoIncrement" json:"id"`
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
	Latitude                *string   `json:"latitude"`
	Longitude               *string   `json:"longitude"`
	Ownership               *string   `json:"ownership"`
	Pic                     *string   `json:"pic"`
	PicInternal             *string   `json:"pic_internal"`
	IsSkulid                *uint     `json:"is_skulid"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`

	SchoolDetail *AccountTypeSchoolDetails `gorm:"foreignKey:AccountID;references:ID"`
}

type AccountTypeSchoolDetails struct {
	AccountID               int        `json:"account_id"`
	DiesNatalis             *time.Time `json:"dies_natalis"`
	Extracurricular         *string    `json:"extracurricular"`
	FootballFieldBrannnding *string    `json:"football_field_brannnding"`
	BasketballFieldBranding *string    `json:"basketball_field_branding"`
	WallPaintingBranding    *string    `json:"wall_painting_branding"`
	WallMagazineBranding    *string    `json:"wall_magazine_branding"`
}

type AccountSimple struct {
	ID              uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	AccountName     *string `json:"account_name"`
	AccountType     *string `json:"account_type"`
	AccountCategory *string `json:"account_category"`
	AccountCode     *string `json:"account_code"`
}
