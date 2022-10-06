package tasks

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/plaid/plaid-go/plaid"
	"github.com/tony-tvu/goexpense/database"
	"github.com/tony-tvu/goexpense/models"
	"github.com/tony-tvu/goexpense/plaidclient"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tasks struct {
	Db           *database.MongoDb
	PlaidClient  *plaidclient.PlaidClient
	TaskInterval int
}

func (t *Tasks) Start(ctx context.Context) {
	go t.refreshTransactionsTask(ctx)
	go t.refreshAccountsTask(ctx)
}

func (t *Tasks) refreshAccountsTask(ctx context.Context) {
	for {
		items, err := database.GetItems(ctx, t.Db)
		if err != nil {
			log.Printf("error getting items: %+v\n", err)
		}

		log.Printf("running accounts scheduled task for %d items. task interval: %ds", len(items), t.TaskInterval)
		for _, item := range items {
			t.RefreshAccountData(ctx, item)
		}
		time.Sleep(time.Duration(t.TaskInterval) * time.Second)
	}
}

func (t *Tasks) refreshTransactionsTask(ctx context.Context) {
	for {
		items, err := database.GetItems(ctx, t.Db)
		if err != nil {
			log.Printf("error getting items: %+v\n", err)
		}

		log.Printf("running transactions scheduled task for %d items. task interval: %ds", len(items), t.TaskInterval)
		for _, item := range items {
			t.RefreshTransactionsData(ctx, item)
		}
		time.Sleep(time.Duration(t.TaskInterval) * time.Second)
	}
}

// function retreives latest accounts data
func (t *Tasks) RefreshAccountData(ctx context.Context, item *models.Item) {
	plaidAccounts, err := t.PlaidClient.GetItemAccounts(ctx, &item.AccessToken)
	if err != nil {
		log.Printf("error getting item's accounts data for plaid_item_id: %s; %+v\n", item.PlaidItemID, err)

		t.Db.Items.UpdateOne(
			ctx,
			bson.M{"plaid_item_id": item.PlaidItemID},
			bson.M{
				"$set": bson.M{
					"item_login_required": true,
					"updated_at":          time.Now()}},
		)
		return
	}

	for _, plaidAccount := range *plaidAccounts {
		count, err := t.Db.Accounts.CountDocuments(ctx, bson.M{"plaid_account_id": plaidAccount.AccountId})
		if err != nil {
			log.Printf("error checking if account exists: %+v\n", err)
		}

		// save new account document if does not exist yet
		if count == 0 {
			doc := &bson.D{
				{Key: "user_id", Value: item.UserID},
				{Key: "plaid_item_id", Value: item.PlaidItemID},
				{Key: "plaid_account_id", Value: plaidAccount.AccountId},
				{Key: "type", Value: plaidAccount.Subtype.Get()},
				{Key: "current_balance", Value: *plaidAccount.Balances.Current.Get()},
				{Key: "name", Value: plaidAccount.Name},
				{Key: "institution", Value: item.Institution},
				{Key: "created_at", Value: time.Now()},
				{Key: "updated_at", Value: time.Now()},
			}
			if _, err := t.Db.Accounts.InsertOne(ctx, doc); err != nil {
				log.Printf("error inserting new account: %+v\n", err)
			}
			continue
		}

		// update existing account document
		_, err = t.Db.Accounts.UpdateOne(
			ctx,
			bson.M{"plaid_account_id": plaidAccount.AccountId},
			bson.M{
				"$set": bson.M{
					"current_balance": *plaidAccount.Balances.Current.Get(),
					"updated_at":      time.Now()}},
		)
		if err != nil {
			log.Printf("error updating item cursor: %+v\n", err)
		}
	}
}

// function retreives latest transactions data
func (t *Tasks) RefreshTransactionsData(ctx context.Context, item *models.Item) {
	isSuccess := true
	transactions, modifiedTransactions, removedTransactions, cursor, err := t.PlaidClient.GetNewTransactions(ctx, item)
	if err != nil {
		log.Printf("error getting transactions for plaid_item_id: %v; err: %+v", item.PlaidItemID, err)
		isSuccess = false
	}

	// save new transactions
	log.Printf("inserting %v transactions for plaid_item_id: %v", len(transactions), item.PlaidItemID)
	for _, transaction := range transactions {
		err := t.saveTransaction(ctx, &transaction, &item.UserID, &item.PlaidItemID)
		if err != nil {
			log.Printf("error inserting new transaction: %+v\n", err)
			isSuccess = false
		}
	}

	// handle modified transactions
	log.Printf("modifying %v transactions for plaid_item_id: %v", len(modifiedTransactions), item.PlaidItemID)
	for _, modified := range modifiedTransactions {
		date, _ := time.Parse("2006-01-02", modified.Date)
		_, err = t.Db.Transactions.UpdateOne(
			ctx,
			bson.M{"transaction_id": modified.GetTransactionId()},
			bson.M{
				"$set": bson.M{
					"amount":     modified.Amount,
					"name":       modified.Name,
					"date":       date,
					"updated_at": time.Now()}},
		)
		if err != nil {
			log.Printf("error updating transaction: %+v\n", err)
			isSuccess = false
		}
	}

	// handle removed transactions
	log.Printf("removing %v transactions for plaid_item_id: %v", len(removedTransactions), item.PlaidItemID)
	for _, modified := range removedTransactions {
		_, err = t.Db.Transactions.DeleteOne(ctx, bson.M{"transaction_id": modified.GetTransactionId()})
		if err != nil {
			log.Printf("error removing transaction: %+v\n", err)
			isSuccess = false
		}
	}

	// Do not save new cursor if there was an error - we want to retry on the next task run
	if isSuccess {
		_, err = t.Db.Items.UpdateOne(
			ctx,
			bson.M{"plaid_item_id": item.PlaidItemID},
			bson.M{
				"$set": bson.M{
					"cursor":     cursor,
					"updated_at": time.Now()}},
		)
		if err != nil {
			log.Printf("error updating item cursor: %+v\n", err)
		}
	}
}

func (t *Tasks) saveTransaction(ctx context.Context, transaction *plaid.Transaction, userID *primitive.ObjectID, plaidItemID *string) error {
	date, _ := time.Parse("2006-01-02", transaction.Date)
	doc := &bson.D{
		{Key: "user_id", Value: userID},
		{Key: "plaid_item_id", Value: plaidItemID},
		{Key: "transaction_id", Value: transaction.GetTransactionId()},
		{Key: "date", Value: date},
		{Key: "amount", Value: transaction.Amount},
		{Key: "category", Value: transaction.Category},
		{Key: "name", Value: transaction.Name},
		{Key: "created_at", Value: time.Now()},
		{Key: "updated_at", Value: time.Now()},
	}
	_, err := t.Db.Transactions.InsertOne(ctx, doc)
	if err != nil && !strings.Contains(err.Error(), "duplicate key error") {
		return err
	}

	return nil
}
