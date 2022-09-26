package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Saves webhooks logging information
type Webhook struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	WebhookType     string             `json:"webhook_type" bson:"webhook_type"`
	WebhookCode     string             `json:"webhook_code" bson:"webhook_code"`
	NewTransactions float32            `json:"new_transactions" bson:"new_transactions"`
	ItemId          string             `json:"item_id" bson:"item_id"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
}
