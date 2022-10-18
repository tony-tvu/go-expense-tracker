package finances

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
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

type Rule struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"`
	Substring string             `json:"substring" bson:"substring"`
	Category  string             `json:"category" bson:"category"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
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

func (h *Handler) CreateRule(c *gin.Context) {
	ctx := c.Request.Context()
	defer c.Request.Body.Close()

	userID, _, err := auth.AuthorizeUser(c, h.Db)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	type Input struct {
		Substring string `json:"substring" validate:"required"`
		Category  string `json:"category" validate:"required"`
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
	if util.ContainsEmpty(input.Substring) {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if !util.Contains(&Categories, input.Category) {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	doc := &bson.D{
		{Key: "user_id", Value: *userID},
		{Key: "substring", Value: input.Substring},
		{Key: "category", Value: input.Category},
		{Key: "created_at", Value: time.Now()},
	}
	_, err = h.Db.Rules.InsertOne(ctx, doc)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// update all transactions with rules
	if !h.applyNewRule(ctx, userID, input.Substring, input.Category) {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) applyNewRule(ctx context.Context, userID *primitive.ObjectID, substring, category string) bool {
	success := true
	var transactions []*Transaction
	cursor, _ := h.Db.Transactions.Find(ctx, bson.M{"user_id": *userID})
	if err := cursor.All(ctx, &transactions); err != nil {
		log.Printf("error updating transaction with new rule: %v", err)
		success = false
	}

	for _, transaction := range transactions {
		if strings.Contains(util.RemoveDuplicateWhitespace(transaction.Name), substring) {
			amount := NormalizeAmount(transaction.Amount, category)

			filter := bson.M{"transaction_id": transaction.TransactionID, "user_id": *userID}
			update := bson.M{"$set": bson.M{"category": category, "amount": amount}}
			_, err := h.Db.Transactions.UpdateOne(ctx, filter, update)
			if err != nil {
				log.Printf("error updating transaction with new rule: %v", err)
				success = false
			}
		}
	}

	return success
}

func (h *Handler) DeleteTransaction(c *gin.Context) {
	ctx := c.Request.Context()
	userID, _, err := auth.AuthorizeUser(c, h.Db)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	transactionID := c.Param("transaction_id")
	if util.ContainsEmpty(transactionID) {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	_, err = h.Db.Transactions.DeleteOne(ctx, bson.M{"transaction_id": transactionID, "user_id": *userID})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteRule(c *gin.Context) {
	ctx := c.Request.Context()
	userID, _, err := auth.AuthorizeUser(c, h.Db)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ruleIDHex := c.Param("rule_id")
	if util.ContainsEmpty(ruleIDHex) {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ruleObjID, err := primitive.ObjectIDFromHex(ruleIDHex)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	_, err = h.Db.Rules.DeleteOne(ctx, bson.M{"_id": ruleObjID, "user_id": *userID})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetRules(c *gin.Context) {
	ctx := c.Request.Context()
	userID, _, err := auth.AuthorizeUser(c, h.Db)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var rules []*Rule
	opts := options.Find().SetSort(bson.D{{Key: "substring", Value: -1}})
	cursor, _ := h.Db.Rules.Find(ctx, bson.M{
		"user_id": *userID,
	}, opts)
	if err = cursor.All(ctx, &rules); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rules": rules,
	})
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
		"user_id": *userID,
	}, opts)
	if err = cursor.All(ctx, &transactions); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	fromDate := time.Now()
	toDate := time.Now()
	if hasFilter {
		fromDate = time.Date(year, util.GetMonth(month), 0, 0, 0, 0, 0, time.UTC)
		toDate = time.Date(year, util.GetMonth(month+1), 1, 0, 0, 0, 0, time.UTC)
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

func (h *Handler) CreateTransaction(c *gin.Context) {
	ctx := c.Request.Context()
	defer c.Request.Body.Close()

	userID, _, err := auth.AuthorizeUser(c, h.Db)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	type Input struct {
		Date     string `json:"date" validate:"required"`
		Name     string `json:"name" validate:"required"`
		Category string `json:"category" validate:"required"`
		Amount   string `json:"amount" validate:"required"`
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

	parsedAmount, err := strconv.ParseFloat(strings.Replace(input.Amount, "-", "", -1), 32)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if parsedAmount == 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	parsed, err := time.Parse(time.RFC1123, input.Date)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// zero out time
	dateZeroed := time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, time.UTC)

	if time.Now().Before(dateZeroed) {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	amount := NormalizeAmount(float32(parsedAmount), input.Category)
	transactionID := uuid.New().String()

	doc := bson.D{
		{Key: "transaction_id", Value: transactionID},
		{Key: "enrollment_id", Value: "user_created"},
		{Key: "name", Value: util.RemoveDuplicateWhitespace(input.Name)},
		{Key: "category", Value: input.Category},
		{Key: "amount", Value: amount},
		{Key: "date", Value: dateZeroed},
		{Key: "user_id", Value: *userID},
		{Key: "account_id", Value: "user_created"},
		{Key: "created_at", Value: time.Now()},
		{Key: "updated_at", Value: time.Now()},
	}

	_, err = h.Db.Transactions.InsertOne(ctx, doc)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
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
		Date          string `json:"date" validate:"required"`
		Name          string `json:"name" validate:"required"`
		Category      string `json:"category" validate:"required"`
		Amount        string `json:"amount" validate:"required"`
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

	parsedAmount, err := strconv.ParseFloat(strings.Replace(input.Amount, "-", "", -1), 32)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if parsedAmount == 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	parsed, err := time.Parse(time.RFC1123, input.Date)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// zero out time
	dateZeroed := time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, time.UTC)

	if time.Now().Before(dateZeroed) {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	amount := NormalizeAmount(float32(parsedAmount), input.Category)
	filter := bson.M{"transaction_id": input.TransactionID, "user_id": *userID}
	update := bson.M{"$set": bson.M{
		"date":     dateZeroed,
		"name":     input.Name,
		"category": input.Category,
		"amount":   amount,
	}}
	_, err = h.Db.Transactions.UpdateOne(ctx, filter, update)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateCategory(c *gin.Context) {
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
		FindOne(ctx, bson.M{"user_id": *userID, "transaction_id": input.TransactionID}).
		Decode(&transaction); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	amount := NormalizeAmount(transaction.Amount, input.Category)

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
	cursor, err := h.Db.Accounts.Find(ctx, bson.M{"user_id": *userID}, opts)
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

// make transaction amount positive if category to changed to 'income'
// make transaction amount negative if category is not 'income'/'ignore'
func NormalizeAmount(amount float32, category string) float32 {
	normalized := amount
	if category == "income" && amount < 0 {
		normalized = -1 * amount
	}
	if category != "income" && category != "ignore" && amount > 0 {
		normalized = -1 * amount
	}
	return normalized
}
