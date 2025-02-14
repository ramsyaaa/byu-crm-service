package models

import "time"

type PerformanceSkulId struct {
	ID             uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	IdSkulid       *string    `gorm:"column:id_skulid" json:"id_skulid"`
	UserId         *uint      `gorm:"column:user_id" json:"user_id"`
	UserType       *string    `gorm:"column:user_type" json:"user_type"`
	RegisteredDate *time.Time `gorm:"column:registered_date" json:"registered_date"`
	Msisdn         *string    `gorm:"column:msisdn" json:"msisdn"`
	Provider       *string    `gorm:"column:provider" json:"provider"`
	AccountId      *uint      `gorm:"column:account_id" json:"account_id"`
	UserName       *string    `gorm:"column:user_name" json:"user_name"`
	FlagNewSales   *int       `gorm:"column:flag_new_sales" json:"flag_new_sales"`
	FlagImei       *int       `gorm:"column:flag_imei" json:"flag_imei"`
	RevMtd         *string    `gorm:"column:rev_mtd" json:"rev_mtd"`
	RevMtdM1       *string    `gorm:"column:rev_mtd_m1" json:"rev_mtd_m1"`
	RevDigital     *string    `gorm:"column:rev_digital" json:"rev_digital"`
	ActivityMtd    *string    `gorm:"column:activity_mtd" json:"activity_mtd"`
	FlagActiveMtd  *string    `gorm:"column:flag_active_mtd" json:"flag_active_mtd"`
	SubdistrictId  *uint      `gorm:"column:subdistrict_id" json:"subdistrict_id"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updated_at"`
}
