package http

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type PerformanceDetail struct {
	ID                 uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Periode            *string    `json:"periode"`
	PeriodeDate        *time.Time `json:"periode_date"`
	EventID            *string    `json:"event_id"`
	PoiID              *string    `json:"poi_id"`
	PoiName            *string    `json:"poi_name"`
	PoiType            *string    `json:"poi_type"`
	EventName          *string    `json:"event_name"`
	EventType          *string    `json:"event_type"`
	EventLocationType  *string    `json:"event_location_type"`
	SalesType          *string    `json:"sales_type"`
	SalesType2         *string    `json:"sales_type_2"`
	CityID             *uint      `json:"city_id"`
	SerialNumberMsisdn *string    `json:"serial_number_msisdn"`
	ScanType           *string    `json:"scan_type"`
	ActiveMsisdn       *string    `json:"active_msisdn"`
	ActiveDate         *time.Time `json:"active_date"`
	ActiveCity         *string    `json:"active_city"`
	Validation         *string    `json:"validation"`
	ValidKpi           bool       `json:"valid_kpi"`
	Revenue            *string    `json:"rev"`
	SaDate             *time.Time `json:"sa_date"`
	SoDate             *time.Time `json:"so_date"`
	NewImei            string     `json:"new_imei"`
	SkulIDDate         *time.Time `json:"skul_id_date"`
	AgentID            *string    `json:"agent_id"`
	UserID             *string    `json:"user_id"`
	UserName           *string    `json:"user_name"`
	UserType           *string    `json:"user_type"`
	UserSubType        *string    `json:"user_sub_type"`
	ScanDate           *time.Time `json:"scan_date"`
	Plan               *string    `json:"plan"`
	TopStatus          bool       `json:"top_status"`
	AccountID          *uint      `json:"account_id"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

type (
	StoreRequest struct {
		Periode            *string    `json:"periode"`
		PeriodeDate        *time.Time `json:"periode_date"`
		EventID            *string    `json:"event_id"`
		PoiID              *string    `json:"poi_id"`
		PoiName            *string    `json:"poi_name"`
		PoiType            *string    `json:"poi_type"`
		EventName          *string    `json:"event_name"`
		EventType          *string    `json:"event_type"`
		EventLocationType  *string    `json:"event_location_type"`
		SalesType          *string    `json:"sales_type"`
		SalesType2         *string    `json:"sales_type_2"`
		CityID             *uint      `json:"city_id"`
		SerialNumberMsisdn *string    `json:"serial_number_msisdn"`
		ScanType           *string    `json:"scan_type"`
		ActiveMsisdn       *string    `json:"active_msisdn"`
		ActiveDate         *time.Time `json:"active_date"`
		ActiveCity         *string    `json:"active_city"`
		Validation         *string    `json:"validation"`
		ValidKpi           bool       `json:"valid_kpi"`
		Revenue            *string    `json:"rev"`
		SaDate             *time.Time `json:"sa_date"`
		SoDate             *time.Time `json:"so_date"`
		NewImei            uint8      `json:"new_imei"`
		SkulIDDate         *time.Time `json:"skul_id_date"`
		AgentID            *string    `json:"agent_id"`
		UserID             *string    `json:"user_id"`
		UserName           *string    `json:"user_name"`
		UserType           *string    `json:"user_type"`
		UserSubType        *string    `json:"user_sub_type"`
		ScanDate           *time.Time `json:"scan_date"`
		Plan               *string    `json:"plan"`
		TopStatus          bool       `json:"top_status"`
		AccountID          *uint      `json:"account_id"`
	}
	StoreResponse PerformanceDetail
)

func (r *StoreRequest) bind(c *fiber.Ctx) error {
	return c.BodyParser(r)
}
