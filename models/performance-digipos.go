package models

import "time"

type PerformanceDigipos struct {
	ID                  uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	IdImport            *string    `gorm:"column:id_import" json:"id_import"`
	TrxType             *string    `gorm:"column:trx_type" json:"trx_type"`
	TransactionId       *string    `gorm:"column:transaction_id" json:"transaction_id"`
	EventId             *string    `gorm:"column:event_id" json:"event_id"`
	Status              *string    `gorm:"column:status" json:"status"`
	StatusDesc          *string    `gorm:"column:status_desc" json:"status_desc"`
	ProductId           *string    `gorm:"column:product_id" json:"product_id"`
	ProductName         *string    `gorm:"column:product_name" json:"product_name"`
	SubProductName      *string    `gorm:"column:sub_product_name" json:"sub_product_name"`
	Price               *string    `gorm:"column:price" json:"price"`
	AdminFee            *string    `gorm:"column:admin_fee" json:"admin_fee"`
	StarPoint           *string    `gorm:"column:star_point" json:"star_point"`
	Msisdn              *string    `gorm:"column:msisdn" json:"msisdn"`
	ProductCategory     *string    `gorm:"column:product_category" json:"product_category"`
	DigiposId           *string    `gorm:"column:digipos_id" json:"digipos_id"`
	EventName           *string    `gorm:"column:event_name" json:"event_name"`
	PaymentMethod       *string    `gorm:"column:payment_method" json:"payment_method"`
	SerialNumber        *string    `gorm:"column:serial_number" json:"serial_number"`
	ClusterId           uint       `gorm:"column:cluster_id" json:"cluster_id"`
	CreatedAt           *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedBy           *string    `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt           *time.Time `gorm:"column:updated_at" json:"updated_at"`
	Code                *string    `gorm:"column:code" json:"code"`
	Name                *string    `gorm:"column:name" json:"name"`
	SalesTerritoryLevel *string    `gorm:"column:sales_territory_level" json:"sales_territory_level"`
	SalesTerritoryValue *string    `gorm:"column:sales_territory_value" json:"sales_territory_value"`
	Wok                 *string    `gorm:"column:wok" json:"wok"`
}

func (PerformanceDigipos) TableName() string {
	return "performances"
}
