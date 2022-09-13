package entity

import (
	"time"
)

/*
Plaid Item:

	An Item represents a login at a financial institution. A single end-user of your application might have accounts at different financial institutions, which means they would have multiple different Items. An Item is not the same as a financial institution account, although every account will be associated with an Item. For example, if a user has one login at their bank that allows them to access both their checking account and their savings account, a single Item would be associated with both of those accounts. Each Item linked within your application will have a corresponding access_token, which is a token that you can use to make API requests related to that specific Item. REF: https://plaid.com/docs/quickstart/glossary/#item
*/
type Item struct {
	ID     uint `gorm:"primarykey"`
	UserID uint
	User   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	ItemID      string 
	AccessToken string
	Cursor      string

	CreatedAt time.Time
	UpdatedAt time.Time
}
