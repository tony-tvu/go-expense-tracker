package entity

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	Username     string
	RefreshToken string
	ExpiresAt    time.Time
}
