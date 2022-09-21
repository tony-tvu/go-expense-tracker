package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/database"
	"github.com/tony-tvu/goexpense/models"
	"github.com/tony-tvu/goexpense/util"
	"go.mongodb.org/mongo-driver/bson"
)

type TransactionHandler struct {
	Db    *database.MongoDb
}

type PageInfoInput struct {
	Page int `json:"page" validate:"required,gte=1"`
}

var PAGE_LIMIT = 50

func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	ctx := c.Request.Context()
	
	userID, _, err := auth.AuthorizeUser(c, h.Db)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// TODO: get page from url instead of body
	var input PageInfoInput
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

	pagination := util.Pagination{
		Limit: PAGE_LIMIT,
		Page:  input.Page,
		Sort:  "date desc",
	}

	var transactions []*models.Transaction
	cursor, err := h.Db.Transactions.Find(ctx, bson.M{"user_id": userID})
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
		"page_info":    pagination,
	})
}
