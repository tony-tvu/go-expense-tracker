package resolvers

import (
	"github.com/plaid/plaid-go/plaid"
	"github.com/tony-tvu/goexpense/graph"
	"gorm.io/gorm"
)

type Resolver struct {
	Db          *gorm.DB
	PlaidClient *plaid.APIClient
}

// Mutation returns graph.MutationResolver implementation.
func (r *Resolver) Mutation() graph.MutationResolver { return &mutationResolver{r} }

// Query returns graph.QueryResolver implementation.
func (r *Resolver) Query() graph.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
