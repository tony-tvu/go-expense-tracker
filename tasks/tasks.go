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
	Db                     *database.MongoDb
	Client                 *plaid.APIClient
	TaskInterval           int
	NewTransactionsChannel chan string
	NewAccountsChannel     chan string
}

func (t *Tasks) Start(ctx context.Context) {
	newTransactionsChan := make(chan string, 3)
	t.NewTransactionsChannel = newTransactionsChan
	newAccountsChan := make(chan string, 3)
	t.NewAccountsChannel = newAccountsChan

	go t.newTransactionsListener(ctx)
	go t.newAccountsListener(ctx)
	go t.refreshTransactionsTask(ctx)
	go t.refreshAccountsTask(ctx)
}

// An item's plaid ID will be sent to NewTransactionsChannel from /api/receive_webhooks whenever
// plaid api sends us a webhook specifying that new transactions are available for an item
// This function handles retrieving the updated transactions for that item.
func (t *Tasks) newTransactionsListener(ctx context.Context) {
	for {
		plaidItemID := <-t.NewTransactionsChannel

		var item *models.Item
		if err := t.Db.Items.FindOne(ctx, bson.M{"plaid_item_id": plaidItemID}).Decode(&item); err != nil {
			log.Printf("error getting new item from db: %v\n", err)
		}

		go t.processNewTransactions(ctx, item)
	}
}

func (t *Tasks) newAccountsListener(ctx context.Context) {
	for {
		newItemIDHex := <-t.NewAccountsChannel

		time.Sleep(10 * time.Second)
		log.Printf("getting initial account data for plaid_item_id: %v\n", newItemIDHex)

		objID, err := primitive.ObjectIDFromHex(newItemIDHex)
		if err != nil {
			log.Printf("error getting new item object id: %v\n", err)
		}

		var item *models.Item
		if err = t.Db.Items.FindOne(ctx, bson.M{"_id": objID}).Decode(&item); err != nil {
			log.Printf("error getting new item from db: %v\n", err)
		}

		go t.refreshAccountData(ctx, item)
	}
}

func (t *Tasks) refreshAccountsTask(ctx context.Context) {
	for {
		items, err := database.GetItems(ctx, t.Db)
		if err != nil {
			log.Printf("error getting items: %+v\n", err)
		}

		log.Printf("running accounts scheduled task for %d items. task interval: %ds", len(items), t.TaskInterval)
		for _, item := range items {
			t.refreshAccountData(ctx, item)
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
			t.processNewTransactions(ctx, item)
		}
		time.Sleep(time.Duration(t.TaskInterval) * time.Second)
	}
}

func (t *Tasks) refreshAccountData(ctx context.Context, item *models.Item) {
	plaidAccounts, err := plaidclient.GetItemAccounts(ctx, &item.AccessToken)
	if err != nil {
		log.Printf("error getting item's accounts: %+v\n", err)
	}

	for _, plaidAccount := range *plaidAccounts {
		if *plaidAccount.Subtype.Get() != "checking" && *plaidAccount.Subtype.Get() != "savings" {
			continue
		}

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

func (t *Tasks) processNewTransactions(ctx context.Context, item *models.Item) {
	isSuccess := true
	transactions, modifiedTransactions, removedTransactions, cursor, err := plaidclient.GetNewTransactions(ctx, item)
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
			bson.M{"_id": item.ID},
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
