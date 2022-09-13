package resolvers

import (
	"context"
	"log"

	"github.com/tony-tvu/goexpense/graph/models"
	"github.com/tony-tvu/goexpense/middleware"
	"github.com/tony-tvu/goexpense/util"
)

func (r *queryResolver) Transactions(ctx context.Context, input *models.TransactionSearchInput) (*models.TransactionConnection, error) {
	c := middleware.GetWriterAndCookies(ctx)

	pagination := util.Pagination{
		Limit: 5,
		Page: 1,
		Sort: "Id desc",
	}

	conn := new(models.TransactionConnection)

	log.Println(c)
	log.Println(conn)

	var transactions []*models.Transaction

	r.Db.Scopes(util.Paginate(transactions, &pagination, r.Db)).Find(&transactions)


	for _, transaction := range transactions {
		log.Printf("%+v", transaction)
	}
	
	return nil, nil
}