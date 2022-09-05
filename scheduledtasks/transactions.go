package scheduledtasks

import (
	"context"
	"log"
	"time"

	"github.com/plaid/plaid-go/plaid"
	"github.com/tony-tvu/goexpense/entity"
)

func (s *ScheduledTasks) StartRefreshTransactionsTask() {
	go func() {
		for {
			log.Println("Running RefreshTransactionsTask")

			var items []entity.Item
			if result := s.Db.Raw("SELECT * FROM items").Scan(&items); result.Error != nil {
				log.Printf("error occurred during refreshTransactions task: %+v\n", result.Error)
			}

			log.Printf("Refreshing transactions for %d items\n", len(items))

			for _, item := range items {
				ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Minute*5))
				defer cancel()

				// New transaction updates since "cursor"
				var newTransactions []plaid.Transaction
				var modified []plaid.Transaction
				var removed []plaid.RemovedTransaction // Removed transaction ids
				hasMore := true

				// Iterate through each page of new transaction updates for item
				for hasMore {
					request := plaid.NewTransactionsSyncRequest(item.AccessToken)
					if item.Cursor != "" {
						request.SetCursor(item.Cursor)
					}
					resp, _, err := s.Client.PlaidApi.TransactionsSync(
						ctx,
					).TransactionsSyncRequest(*request).Execute()
					if err != nil {
						log.Printf("error occurred in RefreshTransactionsTask: %+v\n", err)
						return
					}

					// Add this page of results
					newTransactions = append(newTransactions, resp.GetAdded()...)

					// TODO: handling modified transactions
					modified = append(modified, resp.GetModified()...)

					// TODO: handling removed transactions
					removed = append(removed, resp.GetRemoved()...)
					hasMore = resp.GetHasMore()

					// Update cursor to the next cursor
					item.Cursor = resp.GetNextCursor()
				}

				// TODO: save newTransactions to database
				for _, transaction := range newTransactions {
					date, _ := time.Parse("2006-01-02", transaction.Date)
					if result := s.Db.Create(&entity.Transaction{
						ItemID:        item.ID,
						Item:          item,
						UserID:        item.UserID,
						User:          item.User,
						TransactionID: transaction.GetTransactionId(),
						Date:          date,
						Amount:        transaction.Amount,
						Category:      transaction.Category,
						Name:          transaction.Name,
					}); result.Error != nil {
						log.Printf("error occurred in RefreshTransactionsTask while saving new transaction %+v\n", result.Error)
					}
				}
			}

			time.Sleep(time.Duration(taskInterval) * time.Second)
		}
	}()
}
