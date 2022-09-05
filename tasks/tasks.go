package tasks

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/tony-tvu/goexpense/entity"
	"github.com/tony-tvu/goexpense/plaidapi"
	"gorm.io/gorm"
)

// run go routine every x seconds
var taskInterval int
var db *gorm.DB

func init() {
	godotenv.Load(".env")
	taskIntervalInt, err := strconv.Atoi(os.Getenv("TASK_INTERVAL"))
	if err != nil {
		// default: 10m
		taskInterval = 600
	} else {
		taskInterval = taskIntervalInt
	}
}

func Start(gDb *gorm.DB) {
	db = gDb
	go RefreshTransactions()
}

func RefreshTransactions() {
	for {
		log.Println("Running scheduled task: RefreshTransactions")

		var items []entity.Item
		if result := db.Raw("SELECT * FROM items").Scan(&items); result.Error != nil {
			log.Printf("error occurred during refreshTransactions task: %+v\n", result.Error)
		}

		log.Printf("refreshing transactions for %d items\n", len(items))

		for _, item := range items {
			transactions, _, _, cursor, err := plaidapi.GetTransactions(item)
			if err != nil {
				log.Printf("error occurred while getting transaction for itemID: %s; err: %+v", item.ItemID, err)
				continue
			}

			// save new transactions
			for _, t := range transactions {
				date, _ := time.Parse("2006-01-02", t.Date)
				if result := db.Create(&entity.Transaction{
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
					log.Printf("error occurred in RefreshTransactionsTask while saving new transaction %+v\n", result.Error)
				}
			}

			// TODO: handling modified transactions

			// TODO: handling removed transactions

			// TODO: save new cursor for Item
			db.Exec("UPDATE items SET cursor = ? WHERE id = ?", cursor, item.ID)
		}

		time.Sleep(time.Duration(taskInterval) * time.Second)
	}
}
