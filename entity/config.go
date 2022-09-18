package entity

import (
	"time"
)

type Config struct {
	ID uint `json:"id" gorm:"primarykey"`

	// If false, users will not be able to create accounts from UI or handler routes
	RegistrationEnabled bool `json:"registration_enabled" gorm:"default:false"`

	// Because Plaid charges per request, quotas limits can be set and enabled to
	// control the number of times a user can add a new item to their account
	QuotaEnabled bool `json:"quota_enabled" gorm:"default:true"`
	QuotaLimit   int  `json:"quota_limit" gorm:"default:100"`
	QuotaPerUser int  `json:"quota_per_user" gorm:"default:5"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
