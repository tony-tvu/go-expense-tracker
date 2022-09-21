package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Type string

const (
	AdminUser   Type = "ADMIN"
	RegularUser Type = "REGULAR"
)

type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Username string             `json:"username" bson:"username"`
	Email    string             `json:"email" bson:"email"`
	Password string
	Type     Type `json:"type" bson:"type"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

func GetUserType(s string) Type {
	if s == "ADMIN" {
		return AdminUser
	} else {
		return RegularUser
	}
}
