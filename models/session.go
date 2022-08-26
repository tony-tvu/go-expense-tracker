package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	ObjectID     primitive.ObjectID `bson:"_id" json:"_id"`
	Email        string             `bson:"email"`
	RefreshToken string             `bson:"refresh_token"`
	CreatedAt    time.Time          `bson:"created_at"`
	ExpiresAt    time.Time          `bson:"expires_at"`
}
