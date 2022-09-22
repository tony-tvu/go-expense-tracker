package tasks

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/plaid/plaid-go/plaid"
	"github.com/tony-tvu/goexpense/database"
	"github.com/tony-tvu/goexpense/models"
	"go.mongodb.org/mongo-driver/bson"
)

var taskInterval int
var db *database.MongoDb
var client *plaid.APIClient

func init() {
	godotenv.Load(".env")
	taskIntervalInt, err := strconv.Atoi(os.Getenv("TASK_INTERVAL"))
	if err != nil {
		taskInterval = 600
	} else {
		taskInterval = taskIntervalInt
	}
}

func Start(ctx context.Context, gDb *database.MongoDb, pc *plaid.APIClient) {
	db = gDb
	client = pc
	enabled, err := strconv.ParseBool(os.Getenv("TASKS_ENABLED"))
	if err != nil {
		enabled = false
	}
	if !enabled {
		return
	}

	go RefreshTransactions(ctx)
	go RefreshAccounts(ctx)
}

func RefreshTransactions(ctx context.Context) {

	for {
		items, err := database.GetItems(ctx, db)
		if err != nil {
			log.Printf("error getting items: %+v\n", err)
		}
		log.Printf("refreshing transactions for %d items\n", len(items))

		for _, item := range items {
			isSuccess := true
			transactions, _, _, cursor, err := getTransactions(item)
			if err != nil {
				log.Printf("error getting transaction for item_id: %v; err: %+v", item.ID, err)
				isSuccess = false
			}

			// save new transactions
			for _, t := range transactions {
				date, _ := time.Parse("2006-01-02", t.Date)

				doc := &bson.D{
					{Key: "item_id", Value: item.ID},
					{Key: "user_id", Value: item.UserID},
					{Key: "transaction_id", Value: t.GetTransactionId()},
					{Key: "date", Value: date},
					{Key: "amount", Value: t.Amount},
					{Key: "created_at", Value: time.Now()},
					{Key: "updated_at", Value: time.Now()},
				}
				if _, err = db.Transactions.InsertOne(ctx, doc); err != nil {
					if !strings.Contains(err.Error(), "duplicate key error") {
						log.Printf("error inserting new transaction: %+v\n", err)
						isSuccess = false
					}
				}
			}

			// TODO: handling modified transactions

			// TODO: handling removed transactions

			// Do not save new cursor if there was an error - we want to retry on the next task run
			if isSuccess {
				_, err = db.Items.UpdateOne(
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

		time.Sleep(time.Duration(taskInterval) * time.Second)
	}
}

func RefreshAccounts(ctx context.Context) {
	for {
		items, err := database.GetItems(ctx, db)
		if err != nil {
			log.Printf("error getting items: %+v\n", err)
		}
		log.Printf("refreshing accounts for %d items\n", len(items))

		for _, item := range items {
			plaidAccounts, err := getItemAccounts(ctx, item.AccessToken)
			if err != nil {
				log.Printf("error getting item's plaid accounts: %+v\n", err)
			}

			for _, plaidAccount := range *plaidAccounts {
				if *plaidAccount.Subtype.Get() != "checking" && *plaidAccount.Subtype.Get() != "savings" {
					continue
				}

				count, err := db.Accounts.CountDocuments(ctx, bson.D{
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
					if _, err := db.Accounts.InsertOne(ctx, doc); err != nil {
						log.Printf("error inserting new account: %+v\n", err)
					}
					continue
				}

				// update existing account document
				_, err = db.Accounts.UpdateOne(
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

		time.Sleep(time.Duration(taskInterval) * time.Second)
	}
}

// Returns high-level information about all accounts associated with an item
func getItemAccounts(ctx context.Context, accessToken string) (*[]plaid.AccountBase, error) {
	accountsGetResp, _, err := client.PlaidApi.AccountsGet(ctx).AccountsGetRequest(
		*plaid.NewAccountsGetRequest(accessToken),
	).Execute()
	if err != nil {
		return nil, err
	}

	accounts := accountsGetResp.GetAccounts()
	return &accounts, nil
}

func getTransactions(item *models.Item) ([]plaid.Transaction, []plaid.Transaction, []plaid.RemovedTransaction, string, error) {
	ctx := context.Background()

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
		resp, _, err := client.PlaidApi.TransactionsSync(
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
