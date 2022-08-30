package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Type string

const (
	AdminUser   Type = "Admin"
	RegularUser Type = "Regular"
)

type User struct {
	ObjectID  primitive.ObjectID `json:"_id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"password" bson:"password"`
	Type      Type               `json:"type" bson:"type"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}
