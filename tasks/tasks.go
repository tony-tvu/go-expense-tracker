package tasks

import (
	"context"
	"log"
	"os"
	"strconv"
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

func Start(gDb *database.MongoDb, pc *plaid.APIClient) {
	db = gDb
	client = pc
	enabled, err := strconv.ParseBool(os.Getenv("TASKS_ENABLED"))
	if err != nil {
		enabled = false
	}
	if !enabled {
		return
	}

	go RefreshTransactions()
}

func RefreshTransactions() {
	ctx := context.Background()

	for {
		log.Println("Running scheduled task: RefreshTransactions")

		var items []models.Item
		cursor, err := db.Users.Find(ctx, bson.M{})
		if err != nil {
			log.Printf("error occurred during refreshTransactions task: %+v\n", err)
		}
		if err = cursor.All(ctx, &items); err != nil {
			log.Printf("error occurred during refreshTransactions task: %+v\n", err)
		}

		log.Printf("refreshing transactions for %d items\n", len(items))
		for _, item := range items {
			isSuccess := true
			transactions, _, _, cursor, err := getTransactions(item)
			if err != nil {
				log.Printf("error occurred while getting transaction for item_id: %v; err: %+v", item.ID, err)
				isSuccess = false
			}

			// save new transactions
			for _, t := range transactions {
				date, _ := time.Parse("2006-01-02", t.Date)

				doc := &bson.D{
					{Key: "item_id", Value: item.ID},
					{Key: "user_id", Value: item.UserID},
					{Key: "Transaction_id", Value: t.GetTransactionId()},
					{Key: "date", Value: date},
					{Key: "amount", Value: t.Amount},
					{Key: "created_at", Value: time.Now()},
					{Key: "updated_at", Value: time.Now()},
				}
				if _, err = db.Transactions.InsertOne(ctx, doc); err != nil {
					// TODO: check if error is NOT duplicate contraint error, then set isSuccess = false
					log.Printf("error occurred in RefreshTransactionsTask while saving new transaction %+v\n", err)
					isSuccess = false
				}
			}

			// TODO: handling modified transactions

			// TODO: handling removed transactions

			// TODO: save new cursor for Item

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
					log.Printf("error occurred in RefreshTransactionsTask while saving new transaction %+v\n", err)
				}
			}
		}

		time.Sleep(time.Duration(taskInterval) * time.Second)
	}
}

func getTransactions(item models.Item) ([]plaid.Transaction, []plaid.Transaction, []plaid.RemovedTransaction, string, error) {
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
