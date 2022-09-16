package entity

import (
	"time"

	"github.com/lib/pq"
)

type Transaction struct {
	ID     uint `gorm:"primarykey"`
	ItemID uint
	Item   Item `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID uint
	User   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	TransactionID string `gorm:"unique"`
	Date          time.Time
	Amount        float32
	Category      pq.StringArray `gorm:"type:text[]"`
	Name          string

	CreatedAt time.Time
}
