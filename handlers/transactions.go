package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/entity"
	"gorm.io/gorm"
)

type TransactionHandler struct {
	Db *gorm.DB
}

type PageInfoInput struct {
	Page int `json:"page" validate:"required,gte=1"`
}

var PAGE_LIMIT = 50

func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	id, _, err := auth.VerifyUser(c, h.Db)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

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

	pagination := entity.Pagination{
		Limit: PAGE_LIMIT,
		Page:  input.Page,
		Sort:  "date desc",
	}

	var transactions []*entity.Transaction
	h.Db.Scopes(entity.Paginate(transactions, &pagination, h.Db)).Where("user_id = ?", id).Find(&transactions)

	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
		"page_info":    pagination,
	})
}
