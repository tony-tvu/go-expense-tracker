package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role string

const (
	AdminUser    Role = "Admin"
	ExternalUser Role = "External"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name     string             `json:"name"`
	Email    string             `json:"email"`
	Password string             `json:"password"`
	Role     Role               `json:"role"`
	Verified bool               `json:"verified"`
}
