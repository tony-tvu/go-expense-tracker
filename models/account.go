package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Saves information about item's checking and savings accounts
type Account struct {
	ID     primitive.ObjectID `json:"id" bson:"_id"`
	UserID primitive.ObjectID `json:"user_id" bson:"user_id"`

	PlaidItemID    string  `json:"plaid_item_id" bson:"plaid_item_id"`
	PlaidAccountID string  `json:"plaid_account_id" bson:"plaid_account_id"`
	Type           string  `json:"type" bson:"type"`
	CurrentBalance float64 `json:"current_balance" bson:"current_balance"`
	Name           string  `json:"name" bson:"name"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
