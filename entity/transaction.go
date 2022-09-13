package entity

import (
	"time"

	"github.com/lib/pq"
)

type Transaction struct {
	ID     uint `gorm:"primarykey"`
	ItemID string
	UserID uint

	TransactionID string `gorm:"unique"`
	Date          time.Time
	Amount        float32
	Category      pq.StringArray `gorm:"type:text[]"`
	Name          string

	CreatedAt time.Time
	UpdatedAt time.Time
}
