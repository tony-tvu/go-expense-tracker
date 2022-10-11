package finances

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	UserID       primitive.ObjectID `json:"user_id" bson:"user_id"`
	EnrollmentID string             `json:"enrollment_id" bson:"enrollment_id"`
	AccountID    string             `json:"account_id" bson:"account_id"`

	TransactionID string    `json:"transaction_id" bson:"transaction_id"`
	Category      string    `json:"category" bson:"category"`
	Name          string    `json:"name" bson:"name"`
	Date          time.Time `json:"date" bson:"date"`
	Amount        float32   `json:"amount" bson:"amount"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
