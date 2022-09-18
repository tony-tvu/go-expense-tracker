package entity

import (
	"time"
)

type Config struct {
	ID uint `json:"id" gorm:"primarykey"`

	// If false, users will not be able to create accounts from UI or handler routes
	RegistrationEnabled bool `json:"registration_enabled" gorm:"default:false"`

	// If true, users will have a limited number of times they can add a new item
	QuotaEnabled bool `json:"quota_enabled" gorm:"default:true"`
	QuotaLimit   int  `json:"quota_limit" gorm:"default:5"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
