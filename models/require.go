package models

import "time"

type Require struct {
	ID             uint      `json:"id" gorm:"primary_key"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Amount         float32   `json:"amount" gorm:"not null"`
	UnitMeasure    string    `json:"unit_measure" gorm:"type:varchar(128)"`
	SuggestedPrice float32   `json:"suggested_price"`
	Observation    string    `json:"observation"`

	ProductID     uint `json:"product_id"`
	RequirementID uint `json:"requirement_id"`

	QuotationDetails []QuotationDetail `json:"quotation_details"`
}
