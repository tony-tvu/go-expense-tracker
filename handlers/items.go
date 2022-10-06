package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	. "github.com/gobeam/mongo-go-pagination"
	"github.com/plaid/plaid-go/plaid"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/cache"
	"github.com/tony-tvu/goexpense/database"
	"github.com/tony-tvu/goexpense/models"
	"github.com/tony-tvu/goexpense/plaidclient"
	"github.com/tony-tvu/goexpense/tasks"
	"github.com/tony-tvu/goexpense/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ItemHandler struct {
	Db           *database.MongoDb
	ConfigsCache *cache.Configs
	Tasks        *tasks.Tasks
	WebhooksURL  string
	PlaidClient  *plaidclient.PlaidClient
}

func init() {
	v = validator.New()
}

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

	if _, _, err := auth.AuthorizeUser(c, h.Db); err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	linkToken, err := h.PlaidClient.CreateLinkToken(ctx)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"link_token": linkToken})
}

func (h *ItemHandler) GetUpdateLinkToken(c *gin.Context) {
	ctx := c.Request.Context()
	userObjID, _, err := auth.AuthorizeUser(c, h.Db)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	plaidItemID := c.Param("plaid_item_id")

	// verify item belongs to user
	var item *models.Item
	if err = h.Db.Items.FindOne(ctx, bson.M{"plaid_item_id": plaidItemID}).Decode(&item); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if item.UserID != *userObjID {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	linkToken, err := h.PlaidClient.CreateUpdateLinkToken(ctx, &item.AccessToken)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"update_link_token": linkToken})
}

func (h *ItemHandler) UpdateItemsWebhooksURL(c *gin.Context) {
	ctx := c.Request.Context()

	if _, userType, err := auth.AuthorizeUser(c, h.Db); err != nil || *userType != models.AdminUser {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	type Input struct {
		WebhooksURL string `json:"webhooks_url" validate:"required"`
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

	items, err := database.GetItems(ctx, h.Db)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	for _, item := range items {
		err := h.PlaidClient.UpdateWebhooksURL(ctx, &input.WebhooksURL, &item.AccessToken)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
}

// Returns all items associated with userID
func (h *ItemHandler) GetItems(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, err := auth.AuthorizeUser(c, h.Db)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	configs, err := h.ConfigsCache.GetConfigs()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var items []*models.Item
	p, err := New(h.Db.Items).
		Context(ctx).
		Limit(configs.PageLimit).
		Page(int64(page)).
		Sort("institution", 1).
		Select(bson.M{}).
		Filter(bson.M{"user_id": userID}).
		Decode(&items).Find()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// remove plaid access token - should never expose this info
	for _, item := range items {
		item.AccessToken = ""
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     &items,
		"page_info": p.Pagination,
	})
}

// Adds a new item to user's items collection
// New transactions may not be immedietaly available
// ref: https://plaid.com/docs/api/products/transactions/#transactionssync
func (h *ItemHandler) CreateItem(c *gin.Context) {
	ctx := c.Request.Context()
	userID, _, err := auth.AuthorizeUser(c, h.Db)
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
	accessToken, plaidItemID, institution, err := h.PlaidClient.ExchangePublicToken(ctx, input.PublicToken)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// save new Item
	doc := &bson.D{
		{Key: "institution", Value: *institution},
		{Key: "user_id", Value: *userID},
		{Key: "plaid_item_id", Value: plaidItemID},
		{Key: "access_token", Value: accessToken},
		{Key: "cursor", Value: ""},
		{Key: "new_accounts_available", Value: false},
		{Key: "item_login_required", Value: false},
		{Key: "created_at", Value: time.Now()},
		{Key: "updated_at", Value: time.Now()},
	}
	res, err := h.Db.Items.InsertOne(ctx, doc)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// retreive initial account data
	go func() {
		accountCtx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*15))
		defer cancel()

		itemID := res.InsertedID.(primitive.ObjectID)
		log.Printf("getting initial account data for plaid_item_id: %v\n", itemID.Hex())

		var item *models.Item
		if err = h.Db.Items.FindOne(accountCtx, bson.M{"_id": itemID}).Decode(&item); err != nil {
			log.Printf("error getting new item from db: %v\n", err)
		}

		h.Tasks.RefreshAccountData(accountCtx, item)
	}()
}

// Remove item from user's collection
func (h *ItemHandler) DeleteItem(c *gin.Context) {
	ctx := c.Request.Context()

	userObjID, _, err := auth.AuthorizeUser(c, h.Db)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	plaidItemID := c.Param("plaid_item_id")

	// verify item belongs to user
	var item *models.Item
	if err = h.Db.Items.FindOne(ctx, bson.M{"plaid_item_id": plaidItemID}).Decode(&item); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if item.UserID != *userObjID {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// delete item on plaid
	go h.PlaidClient.RemoveItem(&item.AccessToken)

	// delete accounts associated with item
	_, err = h.Db.Accounts.DeleteMany(ctx, bson.M{"plaid_item_id": plaidItemID})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// delete item
	_, err = h.Db.Items.DeleteOne(ctx, bson.M{"plaid_item_id": plaidItemID})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// delete transactions belonging to item
	_, err = h.Db.Transactions.DeleteMany(ctx, bson.M{"plaid_item_id": plaidItemID})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

// Returns all cash accounts associated with a userID
func (h *ItemHandler) GetAccounts(c *gin.Context) {
	ctx := c.Request.Context()

	userObjID, _, err := auth.AuthorizeUser(c, h.Db)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	opts := options.Find().SetSort(bson.D{{Key: "institution", Value: 1}})
	var accounts []*models.Account
	cursor, err := h.Db.Accounts.Find(ctx, bson.M{"user_id": &userObjID}, opts)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if err = cursor.All(ctx, &accounts); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accounts": accounts,
	})
}

var TRANSACTIONS_WEBHOOKS = []string{
	"SYNC_UPDATES_AVAILABLE",
	"INITIAL_UPDATE",
	"HISTORICAL_UPDATE",
	"DEFAULT_UPDATE",
	"TRANSACTIONS_REMOVED",
}

func (h *ItemHandler) ReceiveWebooks(c *gin.Context) {
	ctx := c.Request.Context()
	defer c.Request.Body.Close()

	signedJwt := c.Request.Header.Get("Plaid-Verification")

	sha256Val, err := h.PlaidClient.VerifyWebhook(ctx, signedJwt)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// verify webhook key and body sha matches
	shaCompute := sha256.New()
	shaCompute.Write(bodyBytes)
	bodySha := shaCompute.Sum(nil)
	bodyShaStr := fmt.Sprintf("%x", bodySha)

	if bodyShaStr != *sha256Val {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var webhook *plaid.DefaultUpdateWebhook
	err = json.Unmarshal(bodyBytes, &webhook)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if util.Contains(&TRANSACTIONS_WEBHOOKS, webhook.WebhookCode) {
		log.Printf("webhook received: %+v", webhook)
		go func() {
			transactionsCtx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*15))
			defer cancel()

			var item *models.Item
			if err := h.Db.Items.FindOne(transactionsCtx, bson.M{"plaid_item_id": webhook.ItemId}).Decode(&item); err != nil {
				log.Printf("error getting new item from db: %v\n", err)
			}

			h.Tasks.RefreshTransactionsData(transactionsCtx, item)
		}()
	}

	if webhook.WebhookCode == "ITEM_LOGIN_REQUIRED" || webhook.WebhookCode == "PENDING_EXPIRATION" {
		_, err = h.Db.Items.UpdateOne(
			ctx,
			bson.M{"plaid_item_id": webhook.ItemId},
			bson.M{
				"$set": bson.M{
					"item_login_required": true,
					"updated_at":          time.Now()}},
		)
		if err != nil {
			log.Printf("error updating item_login_required webhook for plaid_item_id: %s; %+v\n", webhook.ItemId, err)
		}
	}

	if webhook.WebhookCode == "NEW_ACCOUNTS_AVAILABLE" {
		_, err = h.Db.Items.UpdateOne(
			ctx,
			bson.M{"plaid_item_id": webhook.ItemId},
			bson.M{
				"$set": bson.M{
					"new_accounts_available": true,
					"updated_at":             time.Now()}},
		)
		if err != nil {
			log.Printf("error updating new_accounts_available webhook for plaid_item_id: %s; %+v\n", webhook.ItemId, err)
		}
	}
}
