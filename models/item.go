package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
Plaid Item:

	An Item represents a login at a financial institution. A single end-user of your application might have accounts at different financial institutions, which means they would have multiple different Items. An Item is not the same as a financial institution account, although every account will be associated with an Item. For example, if a user has one login at their bank that allows them to access both their checking account and their savings account, a single Item would be associated with both of those accounts. Each Item linked within your application will have a corresponding access_token, which is a token that you can use to make API requests related to that specific Item. REF: https://plaid.com/docs/quickstart/glossary/#item
*/
type Item struct {
	ID     primitive.ObjectID `json:"id" bson:"_id"`
	UserID primitive.ObjectID `json:"user_id" bson:"user_id"`

	Institution          string `json:"institution" bson:"institution"`
	PlaidItemID          string `json:"plaid_item_id" bson:"plaid_item_id"`
	AccessToken          string `json:"access_token,omitempty" bson:"access_token"`
	Cursor               string `json:"cursor" bson:"cursor"`
	NewAccountsAvailable bool   `json:"new_accounts_available" bson:"new_accounts_available"`
	ItemLoginRequired    bool   `json:"item_login_required" bson:"item_login_required"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
