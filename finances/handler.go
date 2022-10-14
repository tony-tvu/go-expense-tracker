package finances

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/db"
	"github.com/tony-tvu/goexpense/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Handler struct {
	Db *db.MongoDb
}

type Account struct {
	ID     primitive.ObjectID `json:"id" bson:"_id"`
	UserID primitive.ObjectID `json:"user_id" bson:"user_id"`

	AccountID    string    `json:"account_id" bson:"account_id"`
	EnrollmentID string    `json:"enrollment_id" bson:"enrollment_id"`
	AccessToken  string    `json:"access_token" bson:"access_token"`
	AccountType  string    `json:"account_type" bson:"account_type"`
	Subtype      string    `json:"subtype" bson:"subtype"`
	Status       string    `json:"status" bson:"status"`
	Name         string    `json:"name" bson:"name"`
	LastFour     string    `json:"last_four" bson:"last_four"`
	Institution  string    `json:"institution" bson:"institution"`
	Balance      float64   `json:"balance" bson:"balance"`
	Currency     string    `json:"currency" bson:"currency"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
}

type Transaction struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	UserID       primitive.ObjectID `json:"user_id" bson:"user_id"`
	EnrollmentID string             `json:"enrollment_id" bson:"enrollment_id"`
	AccountID    string             `json:"account_id" bson:"account_id"`

	TransactionID string    `json:"transaction_id" bson:"transaction_id"`
	Category      string    `json:"category" bson:"category"`
	Name          string    `json:"name" bson:"name"`
	Date          time.Time `json:"date" bson:"date"`
	Amount        float32   `json:"amount" bson:"amount"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

var Categories = []string{
	"bills",
	"entertainment",
	"groceries",
	"ignore",
	"income",
	"restaurant",
	"transportation",
	"vacation",
	"uncategorized",
}

var v *validator.Validate

func init() {
	v = validator.New()
}

func (h *Handler) GetTransactions(c *gin.Context) {
	ctx := c.Request.Context()
	userID, _, err := auth.AuthorizeUser(c, h.Db)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	hasFilter := false
	monthStr := c.Query("month")
	yearStr := c.Query("year")
	month := 0
	year := 0
	if !util.ContainsEmpty(monthStr, yearStr) {
		month, err = strconv.Atoi(monthStr)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		year, err = strconv.Atoi(yearStr)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		hasFilter = true
	}

	var transactions []*Transaction
	opts := options.Find().SetSort(bson.D{{Key: "date", Value: -1}})
	cursor, _ := h.Db.Transactions.Find(ctx, bson.M{
		"user_id": &userID,
	}, opts)
	if err = cursor.All(ctx, &transactions); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	fromDate := time.Now()
	toDate := time.Now()
	if hasFilter {
		fromDate = time.Date(year, util.GetMonth(month), 0, 0, 0, 0, 0, time.UTC)
		toDate = time.Date(year, util.GetMonth(month+1), 0, 0, 0, 0, 0, time.UTC)
	}

	filtered := []*Transaction{}
	years := []int{}
	for _, t := range transactions {
		if hasFilter && t.Date.After(fromDate) && t.Date.Before(toDate) {
			filtered = append(filtered, t)
		}

		year := t.Date.Year()
		if !util.ContainsInt(&years, year) {
			years = append(years, year)
		}
	}

	if hasFilter {
		c.JSON(http.StatusOK, gin.H{
			"transactions": filtered,
			"count":        len(filtered),
			"years":        years,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"transactions": transactions,
			"count":        len(transactions),
			"years":        years,
		})
	}
}

func (h *Handler) UpdateTransaction(c *gin.Context) {
	ctx := c.Request.Context()
	defer c.Request.Body.Close()

	userID, _, err := auth.AuthorizeUser(c, h.Db)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	type Input struct {
		TransactionID string `json:"transaction_id" validate:"required"`
		Category      string `json:"category" validate:"required"`
	}

	var input *Input
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

	if !util.Contains(&Categories, input.Category) {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var update primitive.M
	var transaction *Transaction
	if err = h.Db.Transactions.
		FindOne(ctx, bson.M{"user_id": userID, "transaction_id": input.TransactionID}).
		Decode(&transaction); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	amount := transaction.Amount

	// make transaction amount positive if category to changed to 'income'
	if input.Category == "income" && transaction.Amount < 0 {
		amount = -1 * transaction.Amount
	}
	// make transaction amount negative if category is not 'income'/'ignore'
	if input.Category != "income" && input.Category != "ignore" && transaction.Amount > 0 {
		amount = -1 * transaction.Amount
	}

	filter := bson.M{"transaction_id": input.TransactionID, "user_id": *userID}
	update = bson.M{"$set": bson.M{"category": input.Category, "amount": amount}}
	_, err = h.Db.Transactions.UpdateOne(ctx, filter, update)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetAccounts(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, err := auth.AuthorizeUser(c, h.Db)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	opts := options.Find().SetSort(bson.D{{Key: "name", Value: 1}})
	var accounts []*Account
	cursor, err := h.Db.Accounts.Find(ctx, bson.M{"user_id": &userID}, opts)
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
