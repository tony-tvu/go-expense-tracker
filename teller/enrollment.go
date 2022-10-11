package teller

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Enrollment struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	UserID       primitive.ObjectID `json:"user_id" bson:"user_id"`
	EnrollmentID string             `json:"enrollment_id" bson:"enrollment_id"`
	AccessToken  string             `json:"access_token" bson:"access_token"`
	Institution  string             `json:"institution" bson:"institution"`
	Disconnected bool               `json:"disconnected" bson:"disconnected"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}
