package tasks

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/plaid/plaid-go/plaid"
	"github.com/tony-tvu/goexpense/database"
	"github.com/tony-tvu/goexpense/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tasks struct {
	Db              *database.MongoDb
	Client          *plaid.APIClient
	TaskInterval    int
	TasksEnabled    bool
	NewItemsChannel chan string
}

func (t *Tasks) Start(ctx context.Context) {
	if !t.TasksEnabled {
		return
	}
	newItemsChan := make(chan string, 3)
	t.NewItemsChannel = newItemsChan
	go t.processNewItems(ctx)

	go t.refreshTransactionsAndAccountsTask(ctx)
}

func (t *Tasks) processNewItems(ctx context.Context) {
	for {
		newItemIDHex := <-t.NewItemsChannel
		log.Printf("new item received: %v\n", newItemIDHex)

		objID, err := primitive.ObjectIDFromHex(newItemIDHex)
		if err != nil {
			log.Printf("error getting new item object id: %v\n", err)
		}

		var item *models.Item
		if err = t.Db.Items.FindOne(ctx, bson.D{{Key: "_id", Value: objID}}).Decode(&item); err != nil {
			log.Printf("error getting new item from db: %v\n", err)
		}



		go t.getInitialTransactions(ctx, item)
		t.refreshAccounts(ctx, []*models.Item{item})
	}
}

// Plaid api could take several seconds/minutes to populate initial transaction data upon item creation.
// To make sure we get transaction data ASAP for the user, this function makes 10 calls
// to update transactions every 10 seconds upon initial item creation
func (t *Tasks) getInitialTransactions(ctx context.Context, item *models.Item) {
	count := 0

	for count < 10 {
		transactions, _, _, _, err := t.getTransactions(ctx, item)
		log.Println(len(transactions))
		if err != nil {
			log.Printf("error getting initial transactions for item_id: %v; err: %+v", item.ID, err)
		}

		for _, transaction := range transactions {
			err := t.saveTransaction(ctx, transaction, &item.UserID, &item.ID)
			if err != nil {
				log.Printf("error inserting new transaction: %+v\n", err)
			}
		}
		count++
		time.Sleep(10 * time.Second)
	}
}

func (t *Tasks) refreshTransactionsAndAccountsTask(ctx context.Context) {
	for {
		items, err := database.GetItems(ctx, t.Db)
		if err != nil {
			log.Printf("error getting items: %+v\n", err)
		}

		log.Printf("refreshing transactions and accounts for %d items\n", len(items))
		t.refreshTransactions(ctx, items)
		t.refreshAccounts(ctx, items)

		time.Sleep(time.Duration(t.TaskInterval) * time.Second)
	}
}

func (t *Tasks) refreshTransactions(ctx context.Context, items []*models.Item) {
	for _, item := range items {
		isSuccess := true
		transactions, _, _, cursor, err := t.getTransactions(ctx, item)
		if err != nil {
			log.Printf("error getting transactions for item_id: %v; err: %+v", item.ID, err)
			isSuccess = false
		}

		// save new transactions
		for _, transaction := range transactions {
			err := t.saveTransaction(ctx, transaction, &item.UserID, &item.ID)
			if err != nil {
				log.Printf("error inserting new transaction: %+v\n", err)
				isSuccess = false
			}
		}

		// TODO: handling modified transactions
		// TODO: handling removed transactions

		// Do not save new cursor if there was an error - we want to retry on the next task run
		if isSuccess {
			_, err = t.Db.Items.UpdateOne(
				ctx,
				bson.M{"_id": item.ID},
				bson.D{
					{Key: "$set", Value: bson.D{
						{Key: "cursor", Value: cursor},
						{Key: "updated_at", Value: time.Now()},
					}},
				},
			)
			if err != nil {
				log.Printf("error updating item cursor: %+v\n", err)
			}
		}
	}
}

func (t *Tasks) saveTransaction(ctx context.Context, transaction plaid.Transaction, userID, itemID *primitive.ObjectID) error {
	date, _ := time.Parse("2006-01-02", transaction.Date)
	doc := &bson.D{
		{Key: "item_id", Value: itemID},
		{Key: "user_id", Value: userID},
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

func (t *Tasks) refreshAccounts(ctx context.Context, items []*models.Item) {
	for _, item := range items {
		plaidAccounts, err := t.getItemAccounts(ctx, item.AccessToken)
		if err != nil {
			log.Printf("error getting item's plaid accounts: %+v\n", err)
		}

		for _, plaidAccount := range *plaidAccounts {
			if *plaidAccount.Subtype.Get() != "checking" && *plaidAccount.Subtype.Get() != "savings" {
				continue
			}

			count, err := t.Db.Accounts.CountDocuments(ctx, bson.D{
				{Key: "user_id", Value: item.UserID},
				{Key: "account_id", Value: plaidAccount.AccountId},
				{Key: "item_id", Value: item.ID},
			})
			if err != nil {
				log.Printf("error checking if account exists: %+v\n", err)
			}

			// save new account document if does not exist yet
			if count == 0 {
				doc := &bson.D{
					{Key: "user_id", Value: item.UserID},
					{Key: "item_id", Value: item.ID},
					{Key: "account_id", Value: plaidAccount.AccountId},
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
				bson.D{
					{Key: "user_id", Value: item.UserID},
					{Key: "account_id", Value: plaidAccount.AccountId},
					{Key: "item_id", Value: item.ID},
				},
				bson.D{
					{Key: "$set", Value: bson.D{
						{Key: "current_balance", Value: *plaidAccount.Balances.Current.Get()},
						{Key: "updated_at", Value: time.Now()},
					}},
				},
			)
			if err != nil {
				log.Printf("error updating item cursor: %+v\n", err)
			}
		}
	}
}

// Returns high-level information about all accounts associated with an item
func (t *Tasks) getItemAccounts(ctx context.Context, accessToken string) (*[]plaid.AccountBase, error) {
	accountsGetResp, _, err := t.Client.PlaidApi.AccountsGet(ctx).AccountsGetRequest(
		*plaid.NewAccountsGetRequest(accessToken),
	).Execute()
	if err != nil {
		return nil, err
	}

	accounts := accountsGetResp.GetAccounts()
	return &accounts, nil
}

func (t *Tasks) getTransactions(ctx context.Context, item *models.Item) ([]plaid.Transaction, []plaid.Transaction, []plaid.RemovedTransaction, string, error) {
	// New transaction updates since "cursor"
	var transactions []plaid.Transaction
	var modified []plaid.Transaction
	var removed []plaid.RemovedTransaction
	cursor := item.Cursor
	hasMore := true

	for hasMore {
		request := plaid.NewTransactionsSyncRequest(item.AccessToken)
		if cursor != "" {
			request.SetCursor(cursor)
		}
		resp, _, err := t.Client.PlaidApi.TransactionsSync(
			ctx,
		).TransactionsSyncRequest(*request).Execute()
		if err != nil {
			return nil, nil, nil, "", err
		}

		transactions = append(transactions, resp.GetAdded()...)
		modified = append(modified, resp.GetModified()...)
		removed = append(removed, resp.GetRemoved()...)

		hasMore = resp.GetHasMore()
		cursor = resp.GetNextCursor()
	}
	return transactions, modified, removed, cursor, nil
}
