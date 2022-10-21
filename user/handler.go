package user

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/db"
	"github.com/tony-tvu/goexpense/util"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	Db *db.MongoDb
}

var v *validator.Validate

func init() {
	v = validator.New()
}

func (h *Handler) IsLoggedIn(c *gin.Context) {
	_, err := auth.AuthorizeUser(c, h.Db)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"logged_in": false,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"logged_in": true,
		})
	}
}

func (h *Handler) GetUsers(c *gin.Context) {
	ctx := c.Request.Context()
	if _, err := auth.AuthorizeUser(c, h.Db); err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var users []*User
	cursor, err := h.Db.Users.Find(ctx, bson.M{})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if err = cursor.All(ctx, &users); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *Handler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	defer c.Request.Body.Close()

	type Input struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
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

	// find existing user account
	var u *User
	if err = h.Db.Users.FindOne(ctx, bson.M{"username": input.Username}).Decode(&u); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// verify password
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.Password))
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	// create refresh token
	refreshToken, err := auth.GetEncryptedToken(auth.RefreshToken, u.ID.Hex())
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// delete existing sessions
	_, err = h.Db.Sessions.DeleteMany(ctx, bson.M{"user_id": u.ID})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// save new session
	doc := bson.D{
		{Key: "user_id", Value: u.ID},
		{Key: "username", Value: u.Username},
		{Key: "refresh_token", Value: refreshToken.Value},
		{Key: "expires_at", Value: refreshToken.ExpiresAt},
		{Key: "created_at", Value: time.Now()},
		{Key: "updated_at", Value: time.Now()},
	}
	_, err = h.Db.Sessions.InsertOne(ctx, doc)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// create access token
	accessToken, err := auth.GetEncryptedToken(auth.AccessToken, u.ID.Hex())
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	util.SetCookie(c.Writer, "goexpense_access", accessToken.Value, accessToken.ExpiresAt)
	util.SetCookie(c.Writer, "goexpense_refresh", refreshToken.Value, refreshToken.ExpiresAt)
}

func (h *Handler) Logout(c *gin.Context) {
	ctx := c.Request.Context()

	userID, err := auth.AuthorizeUser(c, h.Db)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	_, err = h.Db.Sessions.DeleteMany(ctx, bson.M{"user_id": userID})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	util.SetCookie(c.Writer, "goexpense_access", "", time.Now())
	util.SetCookie(c.Writer, "goexpense_refresh", "", time.Now())
}

func (h *Handler) GetUserInfo(c *gin.Context) {
	ctx := c.Request.Context()

	userID, err := auth.AuthorizeUser(c, h.Db)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var u *User
	if err = h.Db.Users.FindOne(ctx, bson.M{"_id": userID}).Decode(&u); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username": u.Username,
		"email":    u.Email,
	})
}

func (h *Handler) RegisterUser(c *gin.Context) {
	ctx := c.Request.Context()
	defer c.Request.Body.Close()

	type Input struct {
		Username string `json:"username" validate:"required"`
		Email    string `json:"email" validate:"email"`
		Password string `json:"password" validate:"required"`
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	doc := &bson.D{
		{Key: "username", Value: input.Username},
		{Key: "email", Value: input.Email},
		{Key: "password", Value: string(hash)},
		{Key: "created_at", Value: time.Now()},
		{Key: "updated_at", Value: time.Now()},
	}
	_, err = h.Db.Users.InsertOne(ctx, doc)
	if err != nil && strings.Contains(err.Error(), "duplicate key error") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User with that email or username already exists",
		})
		return
	} else if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}
