package entity

import (
	"time"
)

type Config struct {
	ID uint `json:"id" gorm:"primarykey"`

	// If false, users will not be able to create accounts from UI or handler routes
	RegistrationAllowed bool `json:"registration_allowed" gorm:"default:false"`

	// If enforced, users will have a limited number of times they can add a new item
	QuotaEnforced bool `json:"quota_enforced" gorm:"default:true"`
	QuotaLimit    int  `json:"quota_limit" gorm:"default:5"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
