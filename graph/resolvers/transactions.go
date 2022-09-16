package resolvers

import (
	"context"

	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/graph"
	"github.com/tony-tvu/goexpense/middleware"
	"github.com/tony-tvu/goexpense/util"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

var PAGE_LIMIT = 50

func (r *queryResolver) Transactions(ctx context.Context, input graph.TransactionSearchInput) (*graph.TransactionConnection, error) {
	c := middleware.GetWriterAndCookies(ctx)

	id, _, err := auth.VerifyUser(c, r.Db)
	if err != nil {
		return nil, gqlerror.Errorf("not authorized")
	}

	pagination := util.Pagination{
		Limit: PAGE_LIMIT,
		Page:  input.PageInfo.Page,
		Sort:  "date desc",
	}

	conn := new(graph.TransactionConnection)

	var transactions []*graph.Transaction
	r.Db.Scopes(util.Paginate(transactions, &pagination, r.Db)).Where("user_id = ?", id).Find(&transactions)

	conn.Nodes = transactions
	conn.PageInfo = getPageInfo(&pagination)

	return conn, nil
}
