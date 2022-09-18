package entity

import (
	"time"
)

type Session struct {
	ID           uint      `json:"id" gorm:"primarykey"`
	UserID       uint      `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}
