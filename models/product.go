package models

import "time"

type Product struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	State       bool      `json:"state"`

	Requires []Require `json:"requires"`
}
