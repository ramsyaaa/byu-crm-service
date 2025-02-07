package models

import "time"

type PerformanceNami struct {
	ID                 uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Periode            *string    `gorm:"column:periode" json:"periode"`
	PeriodeDate        *time.Time `gorm:"column:periode_date" json:"periode_date"`
	EventID            *string    `gorm:"column:event_id" json:"event_id"`
	PoiID              *string    `gorm:"column:poi_id" json:"poi_id"`
	PoiName            *string    `gorm:"column:poi_name" json:"poi_name"`
	PoiType            *string    `gorm:"column:poi_type" json:"poi_type"`
	EventName          *string    `gorm:"column:event_name" json:"event_name"`
	EventType          *string    `gorm:"column:event_type" json:"event_type"`
	EventLocationType  *string    `gorm:"column:event_location_type" json:"event_location_type"`
	SalesType          *string    `gorm:"column:sales_type" json:"sales_type"`
	SalesType2         *string    `gorm:"column:sales_type_2" json:"sales_type_2"`
	CityID             uint       `gorm:"column:city_id" json:"city_id"`
	SerialNumberMsisdn *string    `gorm:"column:serial_number_msisdn" json:"serial_number_msisdn"`
	ScanType           *string    `gorm:"column:scan_type" json:"scan_type"`
	ActiveMsisdn       *string    `gorm:"column:active_msisdn" json:"active_msisdn"`
	ActiveDate         *time.Time `gorm:"column:active_date" json:"active_date"`
	ActiveCity         *string    `gorm:"column:active_city" json:"active_city"`
	Validation         *string    `gorm:"column:validation" json:"validation"`
	ValidKpi           bool       `gorm:"column:valid_kpi" json:"valid_kpi"`
	Revenue            *string    `gorm:"column:rev" json:"rev"`
	SaDate             *time.Time `gorm:"column:sa_date" json:"sa_date"`
	SoDate             *time.Time `gorm:"column:so_date" json:"so_date"`
	NewImei            string     `gorm:"column:new_imei" json:"new_imei"`
	SkulIDDate         *time.Time `gorm:"column:skul_id_date" json:"skul_id_date"`
	AgentID            *string    `gorm:"column:agent_id" json:"agent_id"`
	UserID             *string    `gorm:"column:user_id" json:"user_id"`
	UserName           *string    `gorm:"column:user_name" json:"user_name"`
	UserType           *string    `gorm:"column:user_type" json:"user_type"`
	UserSubType        *string    `gorm:"column:user_sub_type" json:"user_sub_type"`
	ScanDate           *time.Time `gorm:"column:scan_date" json:"scan_date"`
	Plan               *string    `gorm:"column:plan" json:"plan"`
	TopStatus          bool       `gorm:"column:top_status" json:"top_status"`
	AccountID          uint       `gorm:"column:account_id" json:"account_id"`
	CreatedAt          time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt          time.Time  `gorm:"column:updated_at" json:"updated_at"`
}
