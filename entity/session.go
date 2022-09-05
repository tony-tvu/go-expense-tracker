package entity

import (
	"time"
)

type Session struct {
	ID           uint      `json:"id" gorm:"primarykey"`
	Username     string    `json:"username"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
