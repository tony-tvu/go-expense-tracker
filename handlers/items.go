package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/plaid/plaid-go/plaid"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/entity"
	"gorm.io/gorm"
)

type ItemHandler struct {
	Db     *gorm.DB
	Client *plaid.APIClient
}

func init() {
	v = validator.New()
}

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
func (h *ItemHandler) GetLinkToken(c *gin.Context) {
	ctx := c.Request.Context()

	if _, _, err := auth.VerifyUser(c, h.Db); err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
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
		h.Client.PlaidApi.LinkTokenCreate(ctx).LinkTokenCreateRequest(*request).Execute()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	linkToken := linkTokenCreateResp.GetLinkToken()

	c.JSON(http.StatusOK, gin.H{"link_token": linkToken})
}

func (h *ItemHandler) GetItems(c *gin.Context) {
	userID, _, err := auth.VerifyUser(c, h.Db)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var items []*entity.Item
	h.Db.Raw("SELECT * FROM items WHERE user_id = ?", userID).Scan(&items)

	// remove plaid item id and access token - should never expose this info
	for _, item := range items {
		item.PlaidItemID = ""
		item.AccessToken = ""
	}

	c.JSON(http.StatusOK, items)
}

func (h *ItemHandler) CreateItem(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, err := auth.VerifyUser(c, h.Db)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	type Input struct {
		PublicToken string `json:"public_token" validate:"required"`
	}

	var input Input
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(bodyBytes, &input)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = v.Struct(input)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// exchange the public_token for a permanent access_token and itemID
	exchangePublicTokenResp, _, err :=
		h.Client.PlaidApi.ItemPublicTokenExchange(ctx).
			ItemPublicTokenExchangeRequest(
				*plaid.NewItemPublicTokenExchangeRequest(input.PublicToken),
			).Execute()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	accessToken := exchangePublicTokenResp.GetAccessToken()
	plaidItemID := exchangePublicTokenResp.GetItemId()

	institution, err := getInstitution(ctx, h.Client, accessToken)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// save new Item
	if result := h.Db.Create(&entity.Item{
		Institution: *institution,
		UserID:      *userID,
		PlaidItemID: plaidItemID,
		AccessToken: accessToken,
	}); result.Error != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func (h *ItemHandler) DeleteItem(c *gin.Context) {
	userID, _, err := auth.VerifyUser(c, h.Db)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	type Input struct {
		ID uint `json:"id" validate:"required"`
	}

	var input Input
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(bodyBytes, &input)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = v.Struct(input)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if result := h.Db.Exec("DELETE FROM items WHERE user_id = ? AND id = ?", userID, input.ID); result.Error != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
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
