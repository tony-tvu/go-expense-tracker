package resolvers

import (
	"context"
	"log"

	// "github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/graph/models"
	// "github.com/tony-tvu/goexpense/middleware"
	"github.com/tony-tvu/goexpense/util"
	// "github.com/vektah/gqlparser/v2/gqlerror"
)

func (r *queryResolver) Transactions(ctx context.Context, input models.TransactionSearchInput) (*models.TransactionConnection, error) {
	// c := middleware.GetWriterAndCookies(ctx)

	// if !auth.IsAuthorized(c, r.Db) {
	// 	return nil, gqlerror.Errorf("not authorized")
	// }

	pagination := util.Pagination{
		Limit: input.PageInfo.Limit,
		Page: input.PageInfo.Page,
		Sort: "date desc",
	}

	conn := new(models.TransactionConnection)


	var transactions []*models.Transaction

	r.Db.Scopes(util.Paginate(transactions, &pagination, r.Db)).Find(&transactions)

	log.Printf("\n%+v\n", pagination)

	conn.Nodes = transactions
	conn.PageInfo = getPageInfo(&pagination)

	return conn, nil
}