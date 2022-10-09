package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Saves information about item's checking and savings accounts
type Account struct {
	ID     primitive.ObjectID `json:"id" bson:"_id"`
	UserID primitive.ObjectID `json:"user_id" bson:"user_id"`

	AccountID    string    `json:"account_id" bson:"account_id"`
	EnrollmentID string    `json:"enrollment_id" bson:"enrollment_id"`
	AccessToken  string    `json:"access_token" bson:"access_token"`
	Type         string    `json:"type" bson:"type"`
	Subtype      string    `json:"subtype" bson:"subtype"`
	Status       string    `json:"status" bson:"status"`
	Name         string    `json:"name" bson:"name"`
	LastFour     string    `json:"last_four" bson:"last_four"`
	Institution  string    `json:"institution" bson:"institution"`
	Balance      float64   `json:"balance" bson:"balance"`
	Currency     string    `json:"currency" bson:"currency"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
}
