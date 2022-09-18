package config

import (
	"database/sql"
	"time"
)

type Config struct {
	ID                uint         `json:"id" gorm:"primarykey"`
	AllowRegistration sql.NullBool `json:"allow_registration" gorm:"default:false"`
	EnforceQuota      sql.NullBool `json:"enforce_quota" gorm:"default:true"`
	QuotaLimit        int          `json:"quota_limit" gorm:"default:5"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
