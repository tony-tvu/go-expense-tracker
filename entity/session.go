package entity

import (
	"time"
)

type Session struct {
	ID           uint `gorm:"primarykey"`
	UserID       uint
	RefreshToken string
	ExpiresAt    time.Time
	CreatedAt    time.Time
}
