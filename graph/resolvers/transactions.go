package resolvers

import (
	"context"
	"log"

	// "github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/graph"
	// "github.com/tony-tvu/goexpense/middleware"
	"github.com/tony-tvu/goexpense/util"
	// "github.com/vektah/gqlparser/v2/gqlerror"
)

func (r *queryResolver) Transactions(ctx context.Context, input graph.TransactionSearchInput) (*graph.TransactionConnection, error) {
	// c := middleware.GetWriterAndCookies(ctx)

	// if !auth.IsAuthorized(c, r.Db) {
	// 	return nil, gqlerror.Errorf("not authorized")
	// }

	pagination := util.Pagination{
		Limit: input.PageInfo.Limit,
		Page: input.PageInfo.Page,
		Sort: "date desc",
	}

	conn := new(graph.TransactionConnection)


	var transactions []*graph.Transaction

	r.Db.Scopes(util.Paginate(transactions, &pagination, r.Db)).Where("user_id = ?", 1).Find(&transactions)

	log.Printf("\n%+v\n", pagination)

	conn.Nodes = transactions
	conn.PageInfo = getPageInfo(&pagination)

	return conn, nil
}