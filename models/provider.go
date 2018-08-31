package models

import "time"

type Provider struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name" gorm:"not null"`
	RUC         string    `json:"ruc" gorm:"type:varchar(15); not null; unique"`
	Manager     string    `json:"manager" gorm:"type:varchar(255)"`
	Email       string    `json:"email" gorm:"type:varchar(64)"`
	Phone       string    `json:"phone" gorm:"type:varchar(32)"`
	Address     string    `json:"address" gorm:"type:varchar(255)"`
	Observation string    `json:"observation"`
	State       bool      `json:"state"`

	Quotations []Quotation `json:"quotations"`
}
