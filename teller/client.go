package teller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/tony-tvu/goexpense/db"
	"github.com/tony-tvu/goexpense/finances"
	"github.com/tony-tvu/goexpense/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TellerClient struct {
	Client *http.Client
	Db     *db.MongoDb
}

var BASE_URL = "https://api.teller.io"

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

type TellerTransactionRes struct {
	TransactionID string `json:"id"`
	AccountID     string `json:"account_id"`
	Type          string `json:"type"`
	Details       struct {
		ProcessingStatus string `json:"processing_status"`
		Category         string `json:"category"`
		Counterparty     struct {
			Type string `json:"type"`
			Name string `json:"name"`
		}
	} `json:"details"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Amount      string `json:"amount"`
	Status      string `json:"status"`
}

// Fetch all accounts for a given access_token from teller api
func (t *TellerClient) FetchAccounts(accessToken *string) (*[]TellerAccountRes, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/accounts", BASE_URL), nil)
	req.SetBasicAuth(*accessToken, "")
	res, err := t.Client.Do(req)
	if err != nil {
		return nil, err
	}

	var tellerAccounts *[]TellerAccountRes
	json.NewDecoder(res.Body).Decode(&tellerAccounts)

	return tellerAccounts, nil
}

// Fetch account balance for a given account_id from teller api
func (t *TellerClient) FetchBalance(account *finances.Account) (*float64, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/accounts/%s/balances", BASE_URL, account.AccountID), nil)
	req.SetBasicAuth(account.AccessToken, "")
	res, err := t.Client.Do(req)
	if err != nil {
		return nil, err
	}

	var tellerBalance *TellerBalanceRes
	json.NewDecoder(res.Body).Decode(&tellerBalance)

	var balanceStr string
	if account.Subtype == "credit_card" {
		balanceStr = tellerBalance.Ledger
	} else {
		balanceStr = tellerBalance.Available
	}
	balance, err := strconv.ParseFloat(balanceStr, 64)
	if err != nil {
		return nil, err
	}

	return &balance, nil
}

// Fetch all transactions for a given account_id from teller api
func (t *TellerClient) FetchTransactions(account *finances.Account) (*[]TellerTransactionRes, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/accounts/%s/transactions", BASE_URL, account.AccountID), nil)
	req.SetBasicAuth(account.AccessToken, "")
	res, err := t.Client.Do(req)
	if err != nil {
		return nil, err
	}

	var tellerTransactions *[]TellerTransactionRes
	json.NewDecoder(res.Body).Decode(&tellerTransactions)

	return tellerTransactions, nil
}

// Fetches and populates initial account information for a given access_token from teller api
func (t *TellerClient) PopulateAccounts(userID *primitive.ObjectID, accessToken, enrollmentID *string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Minute))
	defer cancel()

	retryLimit := 3
	count := 0

	for count != retryLimit {
		success := true
		tellerAccounts, err := t.FetchAccounts(accessToken)
		if err != nil {
			log.Printf("error making teller accounts request for access_token %s: %v", *accessToken, err)
			success = false
		}

		var docs []interface{}
		for _, account := range *tellerAccounts {
			doc := bson.D{
				{Key: "user_id", Value: *userID},
				{Key: "account_id", Value: account.AccountID},
				{Key: "enrollment_id", Value: *enrollmentID},
				{Key: "access_token", Value: *accessToken},
				{Key: "account_type", Value: account.Type},
				{Key: "subtype", Value: account.Subtype},
				{Key: "status", Value: account.Status},
				{Key: "name", Value: account.Name},
				{Key: "institution", Value: account.Institution.Name},
				{Key: "balance", Value: 0},
				{Key: "currency", Value: account.Currency},
				{Key: "last_four", Value: account.LastFour},
				{Key: "created_at", Value: time.Now()},
				{Key: "updated_at", Value: time.Now()},
			}
			docs = append(docs, doc)
		}
		_, err = t.Db.Accounts.InsertMany(ctx, docs, &options.InsertManyOptions{
			Ordered: util.BoolPointer(false),
		})
		if err != nil && !strings.Contains(err.Error(), "duplicate key error") {
			log.Printf("error saving new account for access_token %s: %v", *accessToken, err)
			success = false
		}

		count++
		if success && (len(*tellerAccounts) > 0) {
			count = retryLimit
		}
		if !success {
			time.Sleep(30 * time.Second)
		}
	}

	go t.RefreshBalances(accessToken)
	go t.RefreshTransactions(userID, accessToken)
}

// Updates all account balances for a give access_token
func (t *TellerClient) RefreshBalances(accessToken *string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Minute))
	defer cancel()

	var accounts []*finances.Account
	cursor, _ := t.Db.Accounts.Find(ctx, bson.M{"access_token": *accessToken})
	if err := cursor.All(ctx, &accounts); err != nil {
		log.Printf("error finding accounts for access_token %s: %v", *accessToken, err)
	}

	retryLimit := 3
	count := 0

	for count != retryLimit {
		success := true
		for _, account := range accounts {

			balance, err := t.FetchBalance(account)
			if err != nil {
				log.Printf("error making request for accounts balance for account_id %s: %v", account.AccountID, err)
				success = false
			}

			_, err = t.Db.Accounts.UpdateOne(
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
		}
		if !success {
			if count == retryLimit {
				t.Db.Enrollments.UpdateOne(
					ctx,
					bson.M{"access_token": *accessToken},
					bson.M{
						"$set": bson.M{
							"disconnected": true,
							"updated_at":   time.Now(),
						}},
				)
			}

			time.Sleep(30 * time.Second)
		}
	}
}

// Fetches all transactions for a given access_token and saves them to db
func (t *TellerClient) RefreshTransactions(userID *primitive.ObjectID, accessToken *string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Minute))
	defer cancel()

	var accounts []*finances.Account
	cursor, _ := t.Db.Accounts.Find(ctx, bson.M{"access_token": *accessToken})
	if err := cursor.All(ctx, &accounts); err != nil {
		log.Printf("error finding accounts for access_token %s: %v", *accessToken, err)
	}

	var rules []*finances.Rule
	cursor, _ = t.Db.Rules.Find(ctx, bson.M{"user_id": userID})
	if err := cursor.All(ctx, &rules); err != nil {
		log.Printf("error finding rules for access_token %s: %v", *accessToken, err)
	}

	retryLimit := 3
	count := 0

	for count != retryLimit {
		success := true

		for _, account := range accounts {
			tellerTransactions, err := t.FetchTransactions(account)
			if err != nil {
				log.Printf("error making teller transactions request for access_token %s: %v", *accessToken, err)
				success = false
			}

			// retry if there are no transactions
			if len(*tellerTransactions) == 0 {
				success = false
			}

			var docs []interface{}
			for _, t := range *tellerTransactions {
				if t.Status != "posted" {
					continue
				}
				amount, err := strconv.ParseFloat(t.Amount, 32)
				if err != nil {
					log.Printf("error parsing amount for transaction %v: %v", t, err)
					success = false
				}
				if account.Subtype == "credit_card" {
					amount = -1 * amount
				}

				date, err := time.Parse("2006-01-02", t.Date)
				if err != nil {
					log.Printf("error parsing date for transaction %v: %v", t, err)
					success = false
				}

				category := "uncategorized"
				if util.Contains(&finances.Categories, t.Details.Category) {
					category = t.Details.Category
				}

				// apply rules
				name := t.Description
				for _, rule := range rules {
					if strings.Contains(name, rule.Substring) {

						// make transaction amount positive if category to changed to 'income'
						if rule.Category == "income" && amount < 0 {
							amount = -1 * amount
						}
						// make transaction amount negative if category is not 'income'/'ignore'
						if rule.Category != "income" && rule.Category != "ignore" && amount > 0 {
							amount = -1 * amount
						}

						category = rule.Category
					}
				}

				doc := bson.D{
					{Key: "transaction_id", Value: t.TransactionID},
					{Key: "enrollment_id", Value: account.EnrollmentID},
					{Key: "name", Value: util.RemoveDuplicateWhitespace(name)},
					{Key: "category", Value: category},
					{Key: "amount", Value: amount},
					{Key: "date", Value: date},
					{Key: "user_id", Value: account.UserID},
					{Key: "account_id", Value: account.AccountID},
					{Key: "created_at", Value: time.Now()},
					{Key: "updated_at", Value: time.Now()},
				}
				docs = append(docs, doc)
			}
			_, err = t.Db.Transactions.InsertMany(ctx, docs, &options.InsertManyOptions{
				Ordered: util.BoolPointer(false),
			})
			if err != nil && !strings.Contains(err.Error(), "duplicate key error") {
				log.Printf("error saving transactions for access_token %s: %v", *accessToken, err)
				success = false
			}
		}

		count++
		if success {
			count = retryLimit
		}
		if !success {
			if count == retryLimit {
				t.Db.Enrollments.UpdateOne(
					ctx,
					bson.M{"access_token": *accessToken},
					bson.M{
						"$set": bson.M{
							"disconnected": true,
							"updated_at":   time.Now(),
						}},
				)
			}

			time.Sleep(30 * time.Second)
		}
	}
}

func (t *TellerClient) DeleteAccount(accessToken, accountID *string) error {
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/accounts/%s", BASE_URL, *accountID), nil)
	req.SetBasicAuth(*accessToken, "")
	_, err := t.Client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
