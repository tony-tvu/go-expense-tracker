package entity

import (
	"time"
)

type Config struct {
	ID uint `json:"id" gorm:"primarykey"`

	AccessTokenExp  int `json:"access_token_exp" gorm:"default:900"`
	RefreshTokenExp int `json:"refresh_token_exp" gorm:"default:3600"`

	// Because Plaid charges per request, quotas limits can be set and enabled to
	// control the number of times a user can add a new item to their account
	QuotaEnabled bool `json:"quota_enabled" gorm:"default:true"`
	QuotaLimit   int  `json:"quota_limit" gorm:"default:10"`

	TasksEnabled  bool `json:"tasks_enabled" gorm:"default:true"`
	TasksInterval int  `json:"tasks_interval" gorm:"default:60"`

	// If false, users will not be able to create accounts from UI or handler routes
	RegistrationEnabled bool `json:"registration_enabled" gorm:"default:false"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
