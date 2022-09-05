package scheduledtasks

import (
	"log"
	"time"

	"github.com/tony-tvu/goexpense/entity"
	"github.com/tony-tvu/goexpense/plaidapi"
)

func RefreshTransactions() {
	for {
		log.Println("Running RefreshTransactionsTask")

		var items []entity.Item
		if result := db.Raw("SELECT * FROM items").Scan(&items); result.Error != nil {
			log.Printf("error occurred during refreshTransactions task: %+v\n", result.Error)
		}

		log.Printf("Refreshing transactions for %d items\n", len(items))

		for _, item := range items {
			transactions, _, _, _, err := plaidapi.GetTransactions(item)
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
		}

		time.Sleep(time.Duration(taskInterval) * time.Second)
	}
}
