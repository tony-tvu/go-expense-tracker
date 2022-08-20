package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/plaid/plaid-go/plaid"
	"github.com/tony-tvu/goexpense/app"
)

/*
This endpoint returns a link_token to the client. From the client, use the
link_token to make a request to plaid api (usePlaidLink) which will open an
interface to have the user login to their bank account. On success, plaid api
will return a public_token. Send the public_token back to this api (GetAccessToken)
to make a request to plaid api for a permanent access_token and associated item_id,
which can be used to get the user's transactions.
*/
func CreateLinkToken(ctx context.Context, a *app.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body := make(map[string]string)

		countryCodes := []plaid.CountryCode{}
		for _, countryCodeStr := range strings.Split(a.PlaidClient.CountryCodes, ",") {
			countryCodes = append(countryCodes, plaid.CountryCode(countryCodeStr))
		}
		products := []plaid.Products{}
		for _, productStr := range strings.Split(a.PlaidClient.Products, ",") {
			products = append(products, plaid.Products(productStr))
		}
		user := plaid.LinkTokenCreateRequestUser{
			ClientUserId: time.Now().String(),
		}

		request := plaid.NewLinkTokenCreateRequest(
			"Plaid Quickstart",
			"en",
			countryCodes,
			user,
		)
		request.SetProducts(products)

		linkTokenCreateResp, _, err :=
			a.PlaidClient.ApiClient.PlaidApi.LinkTokenCreate(ctx).LinkTokenCreateRequest(*request).Execute()

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		linkToken := linkTokenCreateResp.GetLinkToken()

		body["link_token"] = linkToken
		jData, err := json.Marshal(body)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		w.Write(jData)
	}
}

func SetAccessToken(ctx context.Context, a *app.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		publicToken := r.Header.Get("Public-Token")
		if publicToken == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest),
				http.StatusBadRequest)
			return
		}

		// exchange the public_token for an access_token
		exchangePublicTokenResp, _, err :=
			a.PlaidClient.ApiClient.PlaidApi.ItemPublicTokenExchange(ctx).
				ItemPublicTokenExchangeRequest(
					*plaid.NewItemPublicTokenExchangeRequest(publicToken),
				).Execute()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest),
				http.StatusBadRequest)
			return
		}

		accessToken := exchangePublicTokenResp.GetAccessToken()
		itemID := exchangePublicTokenResp.GetItemId()

		log.Println(accessToken)
		log.Println(itemID)
	}
}