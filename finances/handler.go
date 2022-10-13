package finances

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/db"
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

func (h *Handler) GetTransactions(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, err := auth.AuthorizeUser(c, h.Db)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	opts := options.Find().SetSort(bson.D{{Key: "date", Value: -1}})
	var transactions []*Transaction
	cursor, err := h.Db.Transactions.Find(ctx, bson.M{"user_id": &userID}, opts)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if err = cursor.All(ctx, &transactions); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
	})
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
