package user

import (
	"encoding/json"
	"log"
	"net/http"
	"net/mail"
	"net/smtp"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/database"
	"github.com/tony-tvu/goexpense/models"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	Db *database.Db
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h UserHandler) Login(c *gin.Context) {
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
	err = h.Db.Users.FindOne(ctx, bson.D{{Key: "email", Value: cred.Email}}).Decode(&u)
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
	refreshToken, err := auth.GetEncryptedRefreshToken(ctx, u)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// delete existing sessions
	_, err = h.Db.Sessions.DeleteOne(ctx, bson.M{"email": u.Email})
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

	_, err = h.Db.Sessions.InsertOne(ctx, doc)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// create access token
	accessToken, err := auth.GetEncryptedAccessToken(ctx, u.Email, string(u.Type))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "goexpense_access",
		Value:    accessToken.Value,
		Expires:  accessToken.ExpiresAt,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

func (h UserHandler) Logout(c *gin.Context) {
	ctx := c.Request.Context()
	cookie, err := c.Request.Cookie("goexpense_access")

	// missing access_token
	if err != nil || cookie.Value == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	_, claims, err := auth.GetClaimsWithValidation(cookie.Value)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	_, err = h.Db.Sessions.DeleteMany(ctx, bson.M{"email": claims.Email})
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func (h UserHandler) GetUserInfo(c *gin.Context) {
	ctx := c.Request.Context()
	cookie, err := c.Request.Cookie("goexpense_access")

	// no access token - make user log in
	if err != nil || cookie.Value == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	_, claims, err := auth.GetClaimsWithValidation(cookie.Value)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var u *models.User
	err = h.Db.Users.FindOne(ctx, bson.D{{Key: "email", Value: claims.Email}}).Decode(&u)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	// do not send back hashed password
	u.Password = ""
	c.JSON(200, u)
}

// Send email invite to new user
func (h UserHandler) InviteUser(c *gin.Context) {
	type Body struct {
		Email string `json:"email"`
	}
	defer c.Request.Body.Close()

	var b Body
	err := json.NewDecoder(c.Request.Body).Decode(&b)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// validate email address
	_, err = mail.ParseAddress(b.Email)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	to := []string{b.Email}

	from := os.Getenv("EMAIL_SENDER")
	password := os.Getenv("EMAIL_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	if from == "" || password == "" || smtpHost == "" || smtpPort == "" {
		log.Println("error - email sender configs are missing")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	auth := smtp.PlainAuth("", from, password, smtpHost)

	subject := "Subject: This is the subject of the mail\n"
	body := "This is the body of the mail"
	message := []byte(subject + body)

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func (h UserHandler) GetSessions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "User sessions called"})
}

// Handler is used to check if user is logged in
func (h UserHandler) Ping(c *gin.Context) {
	c.Writer.WriteHeader(200)
}
