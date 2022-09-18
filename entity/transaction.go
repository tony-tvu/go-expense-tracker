package entity

import (
	"time"

	"github.com/lib/pq"
)

type Transaction struct {
	ID     uint `json:"id" gorm:"primarykey"`
	ItemID uint `json:"item_id"`
	UserID uint `json:"user_id"`

	TransactionID string         `json:"transaction_id" gorm:"unique"`
	Date          time.Time      `json:"date"`
	Amount        float32        `json:"amount"`
	Category      pq.StringArray `json:"category" gorm:"type:text[]"`
	Name          string         `json:"name"`

	CreatedAt time.Time `json:"created_at"`
}
