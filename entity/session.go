package entity

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	Email        string
	RefreshToken string
	ExpiresAt    time.Time
}
