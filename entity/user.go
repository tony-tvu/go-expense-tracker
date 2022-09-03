package entity

import (
	"gorm.io/gorm"
)

type Type string

const (
	AdminUser   Type = "Admin"
	RegularUser Type = "Regular"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
	Type     Type   `json:"type"`
}
