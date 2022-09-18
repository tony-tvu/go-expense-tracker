package entity

import (
	"time"
)

type Type string

const (
	AdminUser   Type = "ADMIN"
	RegularUser Type = "REGULAR"
)

type User struct {
	ID       uint   `json:"id" gorm:"primarykey"`
	Username string `json:"username" gorm:"unique"`
	Email    string `json:"email" gorm:"unique"`
	Password string
	Type     Type `json:"type"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
