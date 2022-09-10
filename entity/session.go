package entity

import (
	"time"
)

type Session struct {
	ID           uint `gorm:"primarykey"`
	Username     string
	RefreshToken string
	ExpiresAt    time.Time
	CreatedAt    time.Time
}
