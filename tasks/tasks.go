package tasks

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/plaid/plaid-go/plaid"
	"github.com/tony-tvu/goexpense/models"
	"gorm.io/gorm"
)

type Tasks struct {
	Db *gorm.DB
	PlaidClient *plaid.APIClient
}

var taskInterval int
var db *gorm.DB
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

func Start(gDb *gorm.DB, pc *plaid.APIClient) {
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
	for {
		log.Println("Running scheduled task: RefreshTransactions")

		var items []models.Item
		if result := db.Raw("SELECT * FROM items").Scan(&items); result.Error != nil {
			log.Printf("error occurred during refreshTransactions task: %+v\n", result.Error)
		}
		log.Printf("refreshing transactions for %d items\n", len(items))

		for _, item := range items {
			isSuccess := true
			transactions, _, _, cursor, err := getTransactions(item)
			if err != nil {
				log.Printf("error occurred while getting transaction for itemID: %s; err: %+v", item.ItemID, err)
				isSuccess = false
			}

			// save new transactions
			for _, t := range transactions {
				date, _ := time.Parse("2006-01-02", t.Date)
				if result := db.Create(&models.Transaction{
					ItemID:        item.ID,
					Item:          item,
					UserID:        item.UserID,
					User:          item.User,
					TransactionID: t.GetTransactionId(),
					Date:          date,
					Amount:        t.Amount,
					Category:      t.Category,
					Name:          t.Name,
				}); result.Error != nil {
					if !strings.Contains(result.Error.Error(), "duplicate key value violates unique constraint") {
						log.Printf("error occurred in RefreshTransactionsTask while saving new transaction %+v\n", result.Error)
						isSuccess = false
					}
				}
			}

			// TODO: handling modified transactions

			// TODO: handling removed transactions

			// TODO: save new cursor for Item

			// Do not save new cursor if there was an error - we want to retry in the next task run
			if isSuccess {
				db.Exec("UPDATE items SET cursor = ? WHERE id = ?", cursor, item.ID)
			}
		}

		time.Sleep(time.Duration(taskInterval) * time.Second)
	}
}
