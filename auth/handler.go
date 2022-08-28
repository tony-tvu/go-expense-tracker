package auth

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/tony-tvu/goexpense/app"
	"github.com/tony-tvu/goexpense/models"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	App *app.App
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h AuthHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	defer c.Request.Body.Close()

	var cred Credentials
	err := json.NewDecoder(c.Request.Body).Decode(&cred)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// validate email address
	_, err = mail.ParseAddress(cred.Email)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// find existing user account
	var u *models.User
	err = h.App.Users.FindOne(ctx, bson.D{{Key: "email", Value: cred.Email}}).Decode(&u)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// verify password
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(cred.Password))
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	// create refresh token
	refreshToken, err := CreateRefreshToken(ctx, h.App, u)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// delete existing sessions
	_, err = h.App.Sessions.DeleteOne(ctx, bson.M{"email": u.Email})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// save new session
	doc := bson.D{
		{Key: "email", Value: u.Email},
		{Key: "refresh_token", Value: refreshToken.Value},
		{Key: "created_at", Value: time.Now()},
		{Key: "expires_at", Value: refreshToken.ExpiresAt},
	}

	_, err = h.App.Sessions.InsertOne(ctx, doc)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// create access token
	accessToken, err := CreateAccessToken(ctx, h.App, u.Email, string(u.Type))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "goexpense_access",
		Value:    accessToken.Value,
		Expires:  accessToken.ExpiresAt,
		HttpOnly: true,
	})
}

func (h AuthHandler) Logout(c *gin.Context) {
	ctx := c.Request.Context()
	cookie, err := c.Request.Cookie("goexpense_access")

	// missing access_token
	if err != nil || cookie.Value == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// decrypt access token
	decrypted, err := Decrypt(h.App.EncryptionKey, cookie.Value)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	tkn, err := jwt.ParseWithClaims(decrypted, &Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(h.App.JwtKey), nil
		})
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	claims := tkn.Claims.(*Claims)

	_, err = h.App.Sessions.DeleteMany(ctx, bson.M{"email": claims.Email})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}
