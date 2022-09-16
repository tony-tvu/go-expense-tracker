package resolvers

import (
	"context"
	"strings"
	"time"

	"github.com/plaid/plaid-go/plaid"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/entity"
	"github.com/tony-tvu/goexpense/graph"
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

	if _, _, err := auth.VerifyUser(c, r.Db); err != nil {
		return "", gqlerror.Errorf("not authorized")
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
		convertCountryCodes(strings.Split(countryCodes, ",")),
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

func (r *mutationResolver) SetAccessToken(ctx context.Context, input graph.PublicTokenInput) (bool, error) {
	c := middleware.GetWriterAndCookies(ctx)

	if _, _, err := auth.VerifyUser(c, r.Db); err != nil {
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

	institution, err := getInstitution(ctx, r.PlaidClient, accessToken)
	if err != nil {
		return false, gqlerror.Errorf("internal server error")
	}

	// save new Item
	claims, _ := auth.ValidateTokenAndGetClaims(c.EncryptedRefreshToken)
	var u *entity.User
	if result := r.Db.Where("id = ?", claims.UserID).First(&u); result.Error != nil {
		return false, gqlerror.Errorf("internal server error")
	}

	if result := r.Db.Create(&entity.Item{
		Institution: *institution,
		UserID:      u.ID,
		ItemID:      itemID,
		AccessToken: accessToken,
	}); result.Error != nil {
		return false, gqlerror.Errorf("internal server error")
	}

	return true, nil
}

func getInstitution(ctx context.Context, client *plaid.APIClient, accessToken string) (*string, error) {
	itemGetResp, _, err := client.PlaidApi.ItemGet(ctx).ItemGetRequest(
		*plaid.NewItemGetRequest(accessToken),
	).Execute()
	if err != nil {
		return nil, err
	}

	institutionGetByIdResp, _, err := client.PlaidApi.InstitutionsGetById(ctx).InstitutionsGetByIdRequest(
		*plaid.NewInstitutionsGetByIdRequest(
			*itemGetResp.GetItem().InstitutionId.Get(),
			convertCountryCodes(strings.Split(countryCodes, ",")),
		),
	).Execute()
	if err != nil {
		return nil, err
	}

	institution := institutionGetByIdResp.GetInstitution().Name
	return &institution, nil
}

func convertCountryCodes(countryCodeStrs []string) []plaid.CountryCode {
	codes := []plaid.CountryCode{}
	for _, countryCodeStr := range countryCodeStrs {
		codes = append(codes, plaid.CountryCode(countryCodeStr))
	}

	return codes
}


func (r *queryResolver) Items(ctx context.Context) ([]*graph.Item, error) {
	c := middleware.GetWriterAndCookies(ctx)

	userID, _, err := auth.VerifyUser(c, r.Db)
	if err != nil {
		return nil, gqlerror.Errorf("not authorized")
	}

	var items []*graph.Item
	r.Db.Raw("SELECT * FROM items WHERE user_id = ?", userID).Scan(&items)

	return items, nil
}

func (r *mutationResolver) DeleteItem(ctx context.Context, input graph.DeleteItemInput) (bool, error) {
	c := middleware.GetWriterAndCookies(ctx)

	userID, _, err := auth.VerifyUser(c, r.Db)
	if err != nil {
		return false, gqlerror.Errorf("not authorized")
	}

	if result := r.Db.Exec("DELETE FROM items WHERE user_id = ? AND id = ?", userID, input.ID); result.Error != nil {
		return false, gqlerror.Errorf("internal server error")
	}

	return true, nil
}