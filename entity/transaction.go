package entity

import (
	"time"

	"github.com/lib/pq"
)

type Transaction struct {
	ID     uint `json:"id,string" gorm:"primarykey"`
	ItemID uint `json:"itemId"`
	Item   Item `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID uint `json:"userId"`
	User   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	TransactionID string         `json:"transactionId" gorm:"unique"`
	Date          time.Time      `json:"date"`
	Amount        float32        `json:"amount"`
	Category      pq.StringArray `json:"category" gorm:"type:text[]"`
	Name          string         `json:"name"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
