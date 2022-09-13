package resolvers

import (
	"context"
	"strings"
	"time"

	"github.com/plaid/plaid-go/plaid"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/entity"
	"github.com/tony-tvu/goexpense/graph/models"
	"github.com/tony-tvu/goexpense/middleware"
	"github.com/tony-tvu/goexpense/util"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

var products string = "transactions"
var countryCodes string = "US,CA"

/*
This resolver returns a link_token to the client. From the client, use the
link_token to make a request to plaid api (usePlaidLink) which will open an
interface to have the user login to their bank account. On success, plaid api
will return a public_token. Send the public_token back to this api (SetAccessToken)
which makes a request with plaid's GetAccessToken and returns a permanent access_token
and associated item_id, which can be used to get the user's transactions.
*/
func (r *queryResolver) LinkToken(ctx context.Context) (string, error) {
	c := middleware.GetWriterAndCookies(ctx)

	if !auth.IsAuthorized(c, r.Db) {
		return "", gqlerror.Errorf("not authorized")
	}

	cc := []plaid.CountryCode{}
	for _, countryCode := range strings.Split(countryCodes, ",") {
		cc = append(cc, plaid.CountryCode(countryCode))
	}
	p := []plaid.Products{}
	for _, product := range strings.Split(products, ",") {
		p = append(p, plaid.Products(product))
	}
	user := plaid.LinkTokenCreateRequestUser{
		ClientUserId: time.Now().String(),
	}

	request := plaid.NewLinkTokenCreateRequest(
		"Plaid Quickstart",
		"en",
		cc,
		user,
	)
	request.SetProducts(p)

	linkTokenCreateResp, _, err :=
		r.PlaidClient.PlaidApi.LinkTokenCreate(ctx).LinkTokenCreateRequest(*request).Execute()

	if err != nil {
		return "", gqlerror.Errorf("internal server error")
	}

	linkToken := linkTokenCreateResp.GetLinkToken()

	return linkToken, nil
}

func (r *mutationResolver) SetAccessToken(ctx context.Context, input models.PublicTokenInput) (bool, error) {
	c := middleware.GetWriterAndCookies(ctx)

	if !auth.IsAuthorized(c, r.Db) {
		return false, gqlerror.Errorf("not authorized")
	}

	if util.ContainsEmpty(input.PublicToken) {
		return false, gqlerror.Errorf("bad request")
	}

	// exchange the public_token for a permanent access_token and itemID
	exchangePublicTokenResp, _, err :=
		r.PlaidClient.PlaidApi.ItemPublicTokenExchange(ctx).
			ItemPublicTokenExchangeRequest(
				*plaid.NewItemPublicTokenExchangeRequest(input.PublicToken),
			).Execute()
	if err != nil {
		return false, gqlerror.Errorf("internal server error")
	}

	accessToken := exchangePublicTokenResp.GetAccessToken()
	itemID := exchangePublicTokenResp.GetItemId()

	// save new Item
	claims, _ := auth.ValidateTokenAndGetClaims(c.EncryptedRefreshToken)
	var u *entity.User
	if result := r.Db.Where("id = ?", claims.UserID).First(&u); result.Error != nil {
		return false, gqlerror.Errorf("internal server error")
	}

	if result := r.Db.Create(&entity.Item{
		UserID:      u.ID,
		ItemID:      itemID,
		AccessToken: accessToken,
	}); result.Error != nil {
		return false, gqlerror.Errorf("internal server error")
	}

	return true, nil
}
