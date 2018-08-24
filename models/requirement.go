package models

import (
	"time"
)

type Requirement struct {
	ID             uint      `json:"id" gorm:"primary_key"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Name           string    `json:"name" gorm:"not null"`
	Place          string    `json:"place" gorm:"type:varchar(128)"`
	Destination    string    `json:"destination" gorm:"type:varchar(128)"`
	EmissionDate   time.Time `json:"emission_date"`
	ExpirationDate time.Time `json:"expiration_date"`
	Claimant       string    `json:"claimant"`
	State          string    `json:"state" gorm:"type:varchar(15)"`

	UserID     uint        `json:"user_id"`
	Requires   []Require   `json:"requires"`
	Quotations []Quotation `json:"quotations"`
}
