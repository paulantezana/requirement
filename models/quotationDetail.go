package models

import "time"

type QuotationDetail struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UnitPrice float32   `json:"unit_price" gorm:"not null"`

	RequireID   uint `json:"require_id"`
	QuotationID uint `json:"quotation_id"`

	WinnerLevelProvider uint `json:"winner_level_provider"` // Winner calculate system  by product
	WinnerProviderID    uint `json:"winner_provider_id"`    // Reference provider whiner by product
}
