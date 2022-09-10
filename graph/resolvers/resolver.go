package resolvers

import (
	"context"

	"github.com/plaid/plaid-go/plaid"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/graph"
	"github.com/tony-tvu/goexpense/middleware"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"gorm.io/gorm"
)

type Resolver struct {
	Db *gorm.DB
	PlaidClient *plaid.APIClient
}

// Used to verify if user is currently logged in
func (r *queryResolver) Ping(ctx context.Context) (bool, error) {
	c := middleware.GetWriterAndCookies(ctx)
	if !auth.IsAuthorized(c, r.Db) {
		return false, gqlerror.Errorf("not authorized")
	}
	
	return true, nil
}

// Mutation returns graph.MutationResolver implementation.
func (r *Resolver) Mutation() graph.MutationResolver { return &mutationResolver{r} }

// Query returns graph.QueryResolver implementation.
func (r *Resolver) Query() graph.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
