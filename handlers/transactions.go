package handlers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/cache"
	"github.com/tony-tvu/goexpense/database"
	"github.com/tony-tvu/goexpense/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TransactionHandler struct {
	Db           *database.MongoDb
	ConfigsCache *cache.Configs
}

func (h *TransactionHandler) GetTransactions(c *gin.Context) {
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
	if page <= 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	configs, err := h.ConfigsCache.GetConfigs()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	filter := bson.M{"user_id": userID}
	findOptions := options.Find().
		SetSort(bson.D{{Key: "date", Value: -1}}).
		SetSkip((int64(page) - 1) * configs.PageLimit).
		SetLimit(configs.PageLimit)

	total, err := h.Db.Transactions.CountDocuments(ctx, filter)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	cursor, err := h.Db.Transactions.Find(ctx, filter, findOptions)
	defer cursor.Close(ctx)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var transactions []*models.Transaction
	for cursor.Next(ctx) {
		var transaction *models.Transaction
		cursor.Decode(&transaction)
		transactions = append(transactions, transaction)
	}

	c.JSON(http.StatusOK, gin.H{
		"transactions": &transactions,
		"total":        total,
		"page":         page,
		"last_page":    math.Ceil(float64(total/configs.PageLimit)) + 1,
		"next":         page + 1,
	})
}
