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
}

type AccountSimple struct {
	ID              uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	AccountName     *string `json:"account_name"`
	AccountType     *string `json:"account_type"`
	AccountCategory *string `json:"account_category"`
	AccountCode     *string `json:"account_code"`
}
