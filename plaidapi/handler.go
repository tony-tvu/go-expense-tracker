package plaidapi

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/plaid/plaid-go/plaid"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/entity"
	"github.com/tony-tvu/goexpense/util"
	"gorm.io/gorm"
)

type PlaidHandler struct {
	Db *gorm.DB
}

var envs = map[string]plaid.Environment{
	"sandbox":     plaid.Sandbox,
	"development": plaid.Development,
	"production":  plaid.Production,
}

var client *plaid.APIClient
var products string
var countryCodes string

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("no .env file found")
	}
	clientID := os.Getenv("PLAID_CLIENT_ID")
	secret := os.Getenv("PLAID_SECRET")
	env := os.Getenv("PLAID_ENV")
	if util.ContainsEmpty(clientID, secret, env) {
		log.Println("plaid env variables are missing")
	}

	products = "auth,transactions"
	countryCodes = "US,CA"
	plaidCfg := plaid.NewConfiguration()
	plaidCfg.AddDefaultHeader("PLAID-CLIENT-ID", clientID)
	plaidCfg.AddDefaultHeader("PLAID-SECRET", secret)
	plaidCfg.UseEnvironment(envs[env])
	plaidClient := plaid.NewAPIClient(plaidCfg)
	client = plaidClient
}

/*
This endpoint returns a link_token to the client. From the client, use the
link_token to make a request to plaid api (usePlaidLink) which will open an
interface to have the user login to their bank account. On success, plaid api
will return a public_token. Send the public_token back to this api (SetAccessToken)
which makes a request with plaid's GetAccessToken and returns a permanent access_token
and associated item_id, which can be used to get the user's transactions.
*/
func (h PlaidHandler) CreateLinkToken(c *gin.Context) {
	ctx := c.Request.Context()

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
		client.PlaidApi.LinkTokenCreate(ctx).LinkTokenCreateRequest(*request).Execute()

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

	// exchange the public_token for a permanent access_token and itemID
	exchangePublicTokenResp, _, err :=
		client.PlaidApi.ItemPublicTokenExchange(ctx).
			ItemPublicTokenExchangeRequest(
				*plaid.NewItemPublicTokenExchangeRequest(publicToken),
			).Execute()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	accessToken := exchangePublicTokenResp.GetAccessToken()
	itemID := exchangePublicTokenResp.GetItemId()

	log.Println(accessToken)
	log.Println(itemID)

	// save new Item
	cookie, _ := c.Request.Cookie("goexpense_refresh")
	claims, _ := auth.ValidateTokenAndGetClaims(cookie.Value)
	var u *entity.User
	if result := h.Db.Where("username = ?", claims.Username).First(&u); result.Error != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if result := h.Db.Create(&entity.Item{
		UserID:      u.ID,
		User:        *u,
		ItemID:      itemID,
		AccessToken: accessToken,
	}); result.Error != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}
