package teller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/db"
	"github.com/tony-tvu/goexpense/finances"
	"github.com/tony-tvu/goexpense/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Handler struct {
	Db           *db.MongoDb
	TellerClient *TellerClient
}

var v *validator.Validate

func init() {
	v = validator.New()
}

func (h *Handler) NewEnrollment(c *gin.Context) {
	ctx := c.Request.Context()
	defer c.Request.Body.Close()

	userID, err := auth.AuthorizeUser(c, h.Db)
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
		{Key: "enrollment_id", Value: input.EnrollmentID},
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

	go h.TellerClient.PopulateAccounts(userID, &input.AccessToken, &input.EnrollmentID)
}

func (h *Handler) GetEnrollments(c *gin.Context) {
	ctx := c.Request.Context()

	userID, err := auth.AuthorizeUser(c, h.Db)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	opts := options.Find().SetSort(bson.D{{Key: "institution", Value: 1}})
	var enrollments []*Enrollment
	cursor, err := h.Db.Enrollments.Find(ctx, bson.M{"user_id": *userID}, opts)
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

func (h *Handler) DeleteEnrollment(c *gin.Context) {
	ctx := c.Request.Context()

	userID, err := auth.AuthorizeUser(c, h.Db)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	enrollmentID := c.Param("enrollment_id")
	if util.ContainsEmpty(enrollmentID) {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// delete accounts
	var accounts []*finances.Account
	cursor, err := h.Db.Enrollments.Find(ctx, bson.M{"enrollment_id": enrollmentID, "user_id": *userID})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if err = cursor.All(ctx, &accounts); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	for _, account := range accounts {
		err := h.TellerClient.DeleteAccount(&account.AccessToken, &account.AccountID)
		if err != nil {
			log.Printf("error making teller account delete request for accound_id %s: %v", account.AccountID, err)
		}
	}

	_, err = h.Db.Accounts.DeleteMany(ctx, bson.M{"enrollment_id": enrollmentID, "user_id": *userID})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// delete enrollment
	_, err = h.Db.Enrollments.DeleteOne(ctx, bson.M{"enrollment_id": enrollmentID, "user_id": *userID})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}
