package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/plaid/plaid-go/plaid"
	"github.com/tony-tvu/goexpense/app"
)

func CreateLinkToken(a *app.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
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
