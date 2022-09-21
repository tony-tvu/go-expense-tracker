package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	ID     primitive.ObjectID `json:"id" bson:"_id"`
	ItemID primitive.ObjectID `json:"item_id" bson:"item_id"`
	UserID primitive.ObjectID `json:"user_id" bson:"user_id"`

	TransactionID string    `json:"transaction_id" bson:"transaction_id"`
	Date          time.Time `json:"date" bson:"date"`
	Amount        float32   `json:"amount" bson:"amount"`
	Category      []string  `json:"category"  bson:"category"`
	Name          string    `json:"name" bson:"name"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
