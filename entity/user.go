package entity

import (
	"time"
)

type Type string

const (
	AdminUser   Type = "Admin"
	RegularUser Type = "Regular"
)

type User struct {
	ID       uint   `gorm:"primarykey"`
	Username string `gorm:"unique"`
	Email    string `gorm:"unique"`
	Password string
	Type     Type

	CreatedAt time.Time
	UpdatedAt time.Time
}
