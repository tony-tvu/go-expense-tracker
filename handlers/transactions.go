package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	. "github.com/gobeam/mongo-go-pagination"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/cache"
	"github.com/tony-tvu/goexpense/database"
	"github.com/tony-tvu/goexpense/models"
	"go.mongodb.org/mongo-driver/bson"
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

	configs, err := h.ConfigsCache.GetConfigs()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var transactions []*models.Transaction
	p, err := New(h.Db.Transactions).
		Context(ctx).
		Limit(configs.PageLimit).
		Page(int64(page)).
		Sort("date", -1).
		Select(bson.D{}).
		Filter(bson.M{"user_id": userID}).
		Decode(&transactions).Find()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transactions": &transactions,
		"page_info":    p.Pagination,
	})
}