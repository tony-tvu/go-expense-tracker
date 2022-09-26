package plaidclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"github.com/plaid/plaid-go/plaid"
	"github.com/tony-tvu/goexpense/models"
	"github.com/tony-tvu/goexpense/util"
)

var envs = map[string]plaid.Environment{
	"sandbox":     plaid.Sandbox,
	"development": plaid.Development,
	"production":  plaid.Production,
}

var client *plaid.APIClient
var webhooksURL string
var products string = "transactions"
var countryCodes string = "US"

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("no .env file found")
	}
	clientID := os.Getenv("PLAID_CLIENT_ID")
	secret := os.Getenv("PLAID_SECRET")
	plaidEnv := os.Getenv("PLAID_ENV")
	if util.ContainsEmpty(clientID, secret, plaidEnv) {
		log.Println("plaid env configs are missing - service will not work")
	}
	plaidCfg := plaid.NewConfiguration()
	plaidCfg.AddDefaultHeader("PLAID-CLIENT-ID", clientID)
	plaidCfg.AddDefaultHeader("PLAID-SECRET", secret)
	plaidCfg.UseEnvironment(envs[plaidEnv])
	pc := plaid.NewAPIClient(plaidCfg)
	client = pc
	webhooksURL = os.Getenv("WEBHOOKS_URL")
}

func CreateLinkToken(ctx context.Context) (*string, error) {
	pp := []plaid.Products{}
	for _, product := range strings.Split(products, ",") {
		pp = append(pp, plaid.Products(product))
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
	request.SetProducts(pp)
	request.SetWebhook(webhooksURL)

	linkTokenCreateResp, _, err :=
		client.PlaidApi.LinkTokenCreate(ctx).LinkTokenCreateRequest(*request).Execute()
	if err != nil {
		return nil, err
	}

	linkToken := linkTokenCreateResp.GetLinkToken()

	return &linkToken, nil
}

func UpdateWebhooksURL(ctx context.Context, newURL, accessToken string) error {
	request := plaid.NewItemWebhookUpdateRequest(accessToken)
	request.Webhook = *plaid.NewNullableString(&newURL)
	_, _, err := client.PlaidApi.ItemWebhookUpdate(ctx).ItemWebhookUpdateRequest(*request).Execute()
	if err != nil {
		return err
	}

	return nil
}

// exchange the public_token for a permanent access_token, itemID, and get institution
func ExchangePublicToken(ctx context.Context, publicToken string) (*string, *string, *string, error) {
	exchangePublicTokenResp, _, err :=
		client.PlaidApi.ItemPublicTokenExchange(ctx).
			ItemPublicTokenExchangeRequest(
				*plaid.NewItemPublicTokenExchangeRequest(publicToken),
			).Execute()
	if err != nil {
		return nil, nil, nil, err
	}

	accessToken := exchangePublicTokenResp.GetAccessToken()
	plaidItemID := exchangePublicTokenResp.GetItemId()
	institution, err := getInstitution(ctx, accessToken)
	if err != nil {
		return nil, nil, nil, err
	}

	return &accessToken, &plaidItemID, institution, nil
}

func getInstitution(ctx context.Context, accessToken string) (*string, error) {
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

// Function verifies if webhook came from Plaid api
func VerifyWebhook(ctx context.Context, signedJwt string)  (*string, error) {
	decodedToken, _, err := new(jwt.Parser).ParseUnverified(signedJwt, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}
	if decodedToken.Header["alg"] != "ES256" {
		return nil, errors.New("error - invalid plaid jwt algorithm")
	}

	currentKeyID := fmt.Sprintf("%v", decodedToken.Header["kid"])
	if currentKeyID == "" {
		return nil, errors.New("error - plaid jwt key id (kid) missing")
	}

	webhookReq := plaid.NewWebhookVerificationKeyGetRequest(currentKeyID)
	keyResponse, _, err := client.PlaidApi.WebhookVerificationKeyGet(ctx).WebhookVerificationKeyGetRequest(*webhookReq).Execute()

	type JwkKeys struct {
		Keys []interface{} `json:"keys"`
	}
	jwkKeys := &JwkKeys{
		Keys: []interface{}{keyResponse.GetKey()},
	}

	keyJSON, _ := json.Marshal(jwkKeys)
	jwks, err := keyfunc.NewJSON(keyJSON)
	if err != nil {
		return nil, err
	}

	type Claims struct {
		IssuedAt       float32 `json:"iat"`
		RequestBodySha string  `json:"request_body_sha256"`
		jwt.RegisteredClaims
	}

	token, err := jwt.ParseWithClaims(signedJwt, &Claims{}, jwks.Keyfunc)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("webhook token is invalid")
	}

	// verify that the webhook is not more than 5 minutes old to help prevent replay attacks.
	claims := token.Claims.(*Claims)
	integ, decim := math.Modf(float64(claims.IssuedAt))
	issuedAtTime := time.Unix(int64(integ), int64(decim*(1e9)))
	fiveMinPassed := issuedAtTime.Add(5 * time.Minute)
	hasExpired := fiveMinPassed.Before(time.Now())

	if hasExpired {
		return nil, errors.New("webhook token is over five minutes old")
	}

	if claims.RequestBodySha == "" {
		return nil, errors.New("webhook request body sha is missing")
	}

	return &claims.RequestBodySha, nil
}

// Returns high-level information about all accounts associated with an item
func GetItemAccounts(ctx context.Context, accessToken string) (*[]plaid.AccountBase, error) {
	accountsGetResp, _, err := client.PlaidApi.AccountsGet(ctx).AccountsGetRequest(
		*plaid.NewAccountsGetRequest(accessToken),
	).Execute()
	if err != nil {
		return nil, err
	}

	accounts := accountsGetResp.GetAccounts()
	return &accounts, nil
}

func GetNewTransactions(ctx context.Context, item *models.Item) ([]plaid.Transaction, []plaid.Transaction, []plaid.RemovedTransaction, string, error) {
	// New transaction updates since "cursor"
	var transactions []plaid.Transaction
	var modified []plaid.Transaction
	var removed []plaid.RemovedTransaction
	cursor := item.Cursor
	hasMore := true

	for hasMore {
		request := plaid.NewTransactionsSyncRequest(item.AccessToken)
		if cursor != "" {
			request.SetCursor(cursor)
		}
		resp, _, err := client.PlaidApi.TransactionsSync(
			ctx,
		).TransactionsSyncRequest(*request).Execute()
		if err != nil {
			return nil, nil, nil, "", err
		}

		transactions = append(transactions, resp.GetAdded()...)
		modified = append(modified, resp.GetModified()...)
		removed = append(removed, resp.GetRemoved()...)

		hasMore = resp.GetHasMore()
		cursor = resp.GetNextCursor()
	}
	return transactions, modified, removed, cursor, nil
}