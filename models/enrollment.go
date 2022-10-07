package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Enrollment struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	UserID       primitive.ObjectID `json:"user_id" bson:"user_id"`
	Institution  string             `json:"institution" bson:"institution"`
	AccessToken  string             `json:"access_token" bson:"access_token"`
	Disconnected bool               `json:"disconnected" bson:"disconnected"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}
