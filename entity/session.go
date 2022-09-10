package entity

import (
	"time"
)

type Session struct {
	ID           uint      `json:"id,string" gorm:"primarykey"`
	Username     string    `json:"username"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
	CreatedAt    time.Time `json:"createdAt"`
}
