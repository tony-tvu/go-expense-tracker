package plaidapi

import (
	"context"
	"sort"

	"github.com/plaid/plaid-go/plaid"
)

// Set cursor to empty to receive all historical updates
func GetTransactions(cursor *string, accessToken string) (*[]plaid.Transaction, *string, error) {
	ctx := context.Background()

	// New transaction updates since "cursor"
	var added []plaid.Transaction
	var modified []plaid.Transaction
	var removed []plaid.RemovedTransaction // Removed transaction ids
	hasMore := true

	// Iterate through each page of new transaction updates for item
	for hasMore {
		request := plaid.NewTransactionsSyncRequest(accessToken)
		if cursor != nil {
			request.SetCursor(*cursor)
		}
		resp, _, err := client.PlaidApi.TransactionsSync(
			ctx,
		).TransactionsSyncRequest(*request).Execute()
		if err != nil {
			return nil, nil, err
		}

		// Add this page of results
		added = append(added, resp.GetAdded()...)
		modified = append(modified, resp.GetModified()...)
		removed = append(removed, resp.GetRemoved()...)
		hasMore = resp.GetHasMore()

		// Update cursor to the next cursor
		nextCursor := resp.GetNextCursor()
		cursor = &nextCursor
	}

	sort.Slice(added, func(i, j int) bool {
		return added[i].GetDate() < added[j].GetDate()
	})
	latestTransactions := added[len(added)-9:]

	return &latestTransactions, cursor, nil
}
