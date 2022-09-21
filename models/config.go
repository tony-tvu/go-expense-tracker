package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Config struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`

	AccessTokenExp  int `json:"access_token_exp" bson:"access_token_exp" validate:"required,numeric"`
	RefreshTokenExp int `json:"refresh_token_exp" bson:"refresh_token_exp" validate:"required,numeric"`

	// Because Plaid charges per request, quotas limits can be set and enabled to
	// control the number of times a user can add a new item to their account
	QuotaEnabled bool `json:"quota_enabled" bson:"quota_enabled" validate:"required"`
	QuotaLimit   int  `json:"quota_limit" bson:"quota_limit" validate:"required,numeric"`

	TasksEnabled  bool `json:"tasks_enabled" bson:"tasks_enabled" validate:"required"`
	TasksInterval int  `json:"tasks_interval" bson:"tasks_interval" validate:"required,numeric"`

	// If false, users will not be able to create accounts from UI or handler routes
	RegistrationEnabled bool `json:"registration_enabled" bson:"registration_enabled" validate:"required"`

	PageLimit int64 `json:"page_limit" bson:"page_limit"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
