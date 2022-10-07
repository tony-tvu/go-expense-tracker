package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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

var BASE_URL = "https://api.teller.io"

func (h *TellerHandler) NewEnrollment(c *gin.Context) {
	ctx := c.Request.Context()
	userID, _, err := auth.AuthorizeUser(c, h.Db)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	type Input struct {
		AccessToken  string `json:"access_token" validate:"required"`
		EnrollmentID string `json:"enrollment_id" validate:"required"`
		Institution  string `json:"institution" validate:"required"`
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
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/accounts", BASE_URL), nil)
	req.SetBasicAuth(*accessToken, "")

	res, err := h.TellerClient.Do(req)
	if err != nil {
		log.Printf("error populating accounts for access_token %s: %v", *accessToken, err)
	}

	fmt.Printf("%+v", res)
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
