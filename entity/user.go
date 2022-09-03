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
	ID        uint      `gorm:"primarykey"`
	Username  string    `json:"username" gorm:"unique"`
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `json:"password"`
	Type      Type      `json:"type"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
