package models

import "time"

type Setting struct {
	ID             uint      `json:"id" gorm:"primary_key"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Company        string    `json:"company"`
	Email          string    `json:"email"`
	Identification string    `json:"identification"`
	Logo           string    `json:"logo"`
	City           string    `json:"city"`
	Item           uint      `json:"item"`
	Quotations     uint      `json:"quotations"`
}
