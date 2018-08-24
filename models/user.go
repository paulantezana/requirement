package models

import (
	"mime/multipart"
	"time"
)

type User struct {
	ID          uint                   `json:"id" gorm:"primary_key"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	DNI         string                 `json:"dni" gorm:" type:varchar(15); unique; not null"`
	FirstName   string                 `json:"first_name" gorm:"type:varchar(128)"`
	LastName    string                 `json:"last_name" gorm:"type:varchar(128)"`
	UserName    string                 `json:"user_name" gorm:"type:varchar(64); unique; not null"`
	Gender      string                 `json:"gender"`
	Password    string                 `json:"password" gorm:"type:varchar(64); not null"`
	OldPassword string                 `json:"old_password" gorm:"-"`
	Email       string                 `json:"email" gorm:"type:varchar(64); unique; not null"`
	Avatar      string                 `json:"avatar"`
	Picture     []multipart.FileHeader `json:"picture" gorm:"-"`
	Profile     string                 `json:"profile" gorm:"type:varchar(64)"`
	Key         string                 `json:"key"`
	State       bool                   `json:"state" gorm:"default:'true'"`

	Requirements []Requirement `json:"requirements"`
	Quotations   []Quotation   `json:"quotations"`
}
