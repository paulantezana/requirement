package models

import "time"

type Quotation struct {
	ID            uint      `json:"id" gorm:"primary_key"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	EmissionDate  time.Time `json:"emission_date"`
	Winner        bool      `json:"winner"`         // Final Winner set by admin
	WinnerLevel   uint      `json:"winner_level"`   // Winner casting calculate system
	SuggestWinner bool      `json:"suggest_winner"` // Winner suggestion by user
	DeliverDate   time.Time `json:"deliver_date"`
	Observation   string    `json:"observation"`

	ProviderID    uint `json:"provider_id"`
	UserID        uint `json:"user_id"`
	RequirementID uint `json:"requirement_id"`

	QuotationDetails []QuotationDetail `json:"quotation_details"`
}
