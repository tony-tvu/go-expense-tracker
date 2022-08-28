package plaidapi

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/plaid/plaid-go/plaid"
	"github.com/tony-tvu/goexpense/app"
)

type PlaidHandler struct {
	App *app.App
}

/*
This endpoint returns a link_token to the client. From the client, use the
link_token to make a request to plaid api (usePlaidLink) which will open an
interface to have the user login to their bank account. On success, plaid api
will return a public_token. Send the public_token back to this api (GetAccessToken)
to make a request to plaid api for a permanent access_token and associated item_id,
which can be used to get the user's transactions.
*/
func (h PlaidHandler) CreateLinkToken(c *gin.Context) {
	ctx := c.Request.Context()

	countryCodes := []plaid.CountryCode{}
	for _, countryCodeStr := range strings.Split(h.App.PlaidCountryCodes, ",") {
		countryCodes = append(countryCodes, plaid.CountryCode(countryCodeStr))
	}
	products := []plaid.Products{}
	for _, productStr := range strings.Split(h.App.PlaidProducts, ",") {
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
		h.App.PlaidClient.PlaidApi.LinkTokenCreate(ctx).LinkTokenCreateRequest(*request).Execute()

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	linkToken := linkTokenCreateResp.GetLinkToken()
	c.JSON(http.StatusOK, gin.H{"link_token": linkToken})
}

func (h PlaidHandler) SetAccessToken(c *gin.Context) {
	ctx := c.Request.Context()
	publicToken := c.Request.Header.Get("Plaid-Public-Token")
	if publicToken == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// exchange the public_token for an access_token
	exchangePublicTokenResp, _, err :=
		h.App.PlaidClient.PlaidApi.ItemPublicTokenExchange(ctx).
			ItemPublicTokenExchangeRequest(
				*plaid.NewItemPublicTokenExchangeRequest(publicToken),
			).Execute()
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
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
