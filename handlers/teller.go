package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/cache"
	"github.com/tony-tvu/goexpense/database"
	"github.com/tony-tvu/goexpense/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TellerHandler struct {
	Db           *database.MongoDb
	ConfigsCache *cache.Configs
	TellerClient *http.Client
}

type TellerAccountRes struct {
	AccountID   string `json:"id"`
	Type        string `json:"type"`
	Subtype     string `json:"subtype"`
	Status      string `json:"status"`
	Name        string `json:"name"`
	Institution struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"institution"`
	Currency string `json:"currency"`
	LastFour string `json:"last_four"`
}

type TellerBalanceRes struct {
	AccountID string `json:"account_id"`
	Ledger    string `json:"ledger"`
	Available string `json:"available"`
	Links     struct {
		Self    string `json:"self"`
		Account string `json:"account"`
	} `json:"links"`
}

var BASE_URL = "https://api.teller.io"

func (h *TellerHandler) NewEnrollment(c *gin.Context) {
	ctx := c.Request.Context()
	userID, _, err := auth.AuthorizeUser(c, h.Db)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	type Input struct {
		AccessToken string `json:"access_token" validate:"required"`
		Institution string `json:"institution" validate:"required"`
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

	doc := &bson.D{
		{Key: "user_id", Value: *userID},
		{Key: "institution", Value: input.Institution},
		{Key: "access_token", Value: input.AccessToken},
		{Key: "disconnected", Value: false},
		{Key: "created_at", Value: time.Now()},
		{Key: "updated_at", Value: time.Now()},
	}
	_, err = h.Db.Enrollments.InsertOne(ctx, doc)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	go h.populateAccounts(userID, &input.AccessToken)
}

func (h *TellerHandler) populateAccounts(userID *primitive.ObjectID, accessToken *string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Minute))
	defer cancel()

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/accounts", BASE_URL), nil)
	req.SetBasicAuth(*accessToken, "")

	retryLimit := 3
	count := 0

	for count != retryLimit {
		success := true
		res, err := h.TellerClient.Do(req)
		if err != nil {
			log.Printf("error making teller accounts request for access_token %s: %v", *accessToken, err)
			success = false
		}

		var tellerAccounts *[]TellerAccountRes
		json.NewDecoder(res.Body).Decode(&tellerAccounts)

		for _, account := range *tellerAccounts {
			doc := &bson.D{
				{Key: "user_id", Value: *userID},
				{Key: "account_id", Value: account.AccountID},
				{Key: "access_token", Value: *accessToken},
				{Key: "type", Value: account.Type},
				{Key: "subtype", Value: account.Subtype},
				{Key: "status", Value: account.Status},
				{Key: "name", Value: account.Name},
				{Key: "institution", Value: account.Institution.Name},
				{Key: "currency", Value: account.Currency},
				{Key: "last_four", Value: account.LastFour},
				{Key: "created_at", Value: time.Now()},
				{Key: "updated_at", Value: time.Now()},
			}
			_, err = h.Db.Accounts.InsertOne(ctx, doc)
			if err != nil {
				log.Printf("error saving new account for access_token %s: %v", *accessToken, err)
				success = false
			}
		}

		count++
		if success && (len(*tellerAccounts) > 0) {
			count = retryLimit
		}
		if !success {
			time.Sleep(30 * time.Second)
		}
	}

	go h.RefreshAccountBalance(accessToken)
}

func (h *TellerHandler) RefreshAccountBalance(accessToken *string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Minute))
	defer cancel()

	var accounts []*models.Account
	cursor, _ := h.Db.Accounts.Find(ctx, bson.M{"access_token": *accessToken})
	if err := cursor.All(ctx, &accounts); err != nil {
		log.Printf("error finding accounts for access_token %s: %v", *accessToken, err)
	}

	retryLimit := 3
	count := 0

	for count != retryLimit {
		success := true
		for _, account := range accounts {
			req, _ := http.NewRequest("GET", fmt.Sprintf("%s/accounts/%s/balances", BASE_URL, account.AccountID), nil)
			req.SetBasicAuth(*accessToken, "")

			res, err := h.TellerClient.Do(req)
			if err != nil {
				log.Printf("error making request for accounts balance for account_id %s: %v", account.AccountID, err)
				success = false
			}

			var tellerBalance *TellerBalanceRes
			json.NewDecoder(res.Body).Decode(&tellerBalance)

			var balanceStr string
			if account.Subtype == "credit_card" {
				balanceStr = tellerBalance.Ledger
			} else {
				balanceStr = tellerBalance.Available
			}
			balance, _ := strconv.ParseFloat(balanceStr, 64)

			_, err = h.Db.Accounts.UpdateOne(
				ctx,
				bson.M{"account_id": account.AccountID},
				bson.M{
					"$set": bson.M{
						"balance":    balance,
						"updated_at": time.Now(),
					}},
			)
			if err != nil {
				log.Printf("error updating account balance for account_id %s: %v", account.AccountID, err)
				success = false
			}
		}

		count++
		if success {
			count = retryLimit
		} else {
			time.Sleep(30 * time.Second)
		}
	}
}

func (h *TellerHandler) GetEnrollments(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, err := auth.AuthorizeUser(c, h.Db)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	opts := options.Find().SetSort(bson.D{{Key: "institution", Value: 1}})
	var enrollments []*models.Enrollment
	cursor, err := h.Db.Enrollments.Find(ctx, bson.M{"user_id": &userID}, opts)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if err = cursor.All(ctx, &enrollments); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"enrollments": enrollments,
	})
}

func (h *TellerHandler) GetTransactions(c *gin.Context) {

	req, _ := http.NewRequest("GET", "https://api.teller.io/accounts", nil)
	req.SetBasicAuth("ACCESS_TOKEN", "")

	r, err := h.TellerClient.Do(req)

	if err != nil {
		fmt.Print(err)
	}

	type Account struct {
		Type    string `json:"type"`
		SubType string `json:"subtype"`
		Status  string `json:"status"`
		Name    string `json:"name"`
	}

	var resBody *[]Account
	json.NewDecoder(r.Body).Decode(&resBody)

	c.JSON(http.StatusOK, gin.H{
		"message": resBody,
	})
}
