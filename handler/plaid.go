package handler

import (
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
func CreateLinkToken(a *app.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		body := make(map[string]string)

		countryCodes := []plaid.CountryCode{}
		for _, countryCodeStr := range strings.Split(a.PlaidCountryCodes, ",") {
			countryCodes = append(countryCodes, plaid.CountryCode(countryCodeStr))
		}
		products := []plaid.Products{}
		for _, productStr := range strings.Split(a.PlaidProducts, ",") {
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
			a.PlaidClient.PlaidApi.LinkTokenCreate(ctx).LinkTokenCreateRequest(*request).Execute()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		linkToken := linkTokenCreateResp.GetLinkToken()

		body["link_token"] = linkToken
		jData, err := json.Marshal(body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(jData)
	}
}

func SetAccessToken(a *app.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		publicToken := r.Header.Get("Public-Token")
		if publicToken == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// exchange the public_token for an access_token
		exchangePublicTokenResp, _, err :=
			a.PlaidClient.PlaidApi.ItemPublicTokenExchange(ctx).
				ItemPublicTokenExchangeRequest(
					*plaid.NewItemPublicTokenExchangeRequest(publicToken),
				).Execute()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		accessToken := exchangePublicTokenResp.GetAccessToken()
		itemID := exchangePublicTokenResp.GetItemId()

		log.Println(accessToken)
		log.Println(itemID)

		// TODO: set up user login auth - 
		// 1. get UserID from request
		// 2. persist accessToken, itemID with UserID as foreign key
	}
}
