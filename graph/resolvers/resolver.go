package resolvers

import (
	"github.com/plaid/plaid-go/plaid"
	"github.com/tony-tvu/goexpense/graph"
	"github.com/tony-tvu/goexpense/graph/models"
	"github.com/tony-tvu/goexpense/util"
	"gorm.io/gorm"
)

type Resolver struct {
	Db          *gorm.DB
	PlaidClient *plaid.APIClient
}

func getPageInfo(p *util.Pagination) *models.PageInfo {
	pageInfo := models.PageInfo{
		Limit: p.Limit,
		Page: p.Page,
		Sort: p.Sort,
		TotalRows: int(p.TotalRows),
		TotalPages: p.TotalPages,
	}

	return &pageInfo
}

// Mutation returns graph.MutationResolver implementation.
func (r *Resolver) Mutation() graph.MutationResolver { return &mutationResolver{r} }

// Query returns graph.QueryResolver implementation.
func (r *Resolver) Query() graph.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
