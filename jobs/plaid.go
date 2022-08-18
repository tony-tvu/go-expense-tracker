package jobs

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/plaid/plaid-go/plaid"
)

// Plaid Terminology:
// Item = a connection between the user and their bank account

type PlaidClient struct {
	PlaidClientID     string
	PlaidSecret       string
	PlaidEnv          string
	PlaidProducts     string
	PlaidCountryCodes string
	Client            *plaid.APIClient
}

var environments = map[string]plaid.Environment{
	"sandbox":     plaid.Sandbox,
	"development": plaid.Development,
	"production":  plaid.Production,
}

func (p *PlaidClient) Init(clientID, secret, env, products, countryCodes string) {
	p.PlaidClientID = clientID
	p.PlaidSecret = secret
	p.PlaidEnv = env
	p.PlaidProducts = products
	p.PlaidCountryCodes = countryCodes
	if p.PlaidClientID == "" || p.PlaidSecret == "" {
		log.Println("Error: PLAID_SECRET or PLAID_CLIENT_ID is not set.")
		return
	}

	configuration := plaid.NewConfiguration()
	configuration.AddDefaultHeader("PLAID-CLIENT-ID", p.PlaidClientID)
	configuration.AddDefaultHeader("PLAID-SECRET", p.PlaidSecret)
	configuration.UseEnvironment(environments[p.PlaidEnv])
	client := plaid.NewAPIClient(configuration)
	p.Client = client

	linkToken, err := getLinkToken(p.Client, p.PlaidCountryCodes, p.PlaidProducts)
	if err != nil {
		log.Println("Error creating link token")
		return
	}
	log.Println(linkToken)
}

// Returns a link_token to the client
// From the client, use the link_token to make a request to plaid api which will
// open an interface to have the user login to their bank account. On success, plaid
// api will return a public_token. Use the public_token to make a request to plaid api
// for a permanent access_token and associated item_id (used to get transactions).
func getLinkToken(client *plaid.APIClient, plaidCountryCodes, plaidProducts string) (string, error) {
	ctx := context.Background()
	countryCodes := []plaid.CountryCode{}
	for _, countryCodeStr := range strings.Split(plaidCountryCodes, ",") {
		countryCodes = append(countryCodes, plaid.CountryCode(countryCodeStr))
	}
	products := []plaid.Products{}
	for _, productStr := range strings.Split(plaidProducts, ",") {
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
		client.PlaidApi.LinkTokenCreate(ctx).LinkTokenCreateRequest(*request).Execute()
	if err != nil {
		return "", err
	}
	return linkTokenCreateResp.GetLinkToken(), nil
}
