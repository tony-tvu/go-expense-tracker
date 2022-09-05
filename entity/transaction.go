package entity

import (
	"time"

	"github.com/lib/pq"
)

type Transaction struct {
	ID     uint `json:"id" gorm:"primarykey"`
	ItemID uint `json:"item_id"`
	Item   Item `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID uint `json:"user_id"`
	User   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	TransactionID string         `json:"transaction_id" gorm:"unique"`
	Date          time.Time      `json:"date"`
	Amount        float32        `json:"amount"`
	Category      pq.StringArray `json:"category" gorm:"type:text[]"`
	Name          string         `json:"name"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
