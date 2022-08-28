package user

import (
	"encoding/json"
	"log"
	"net/http"
	"net/mail"
	"net/smtp"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"github.com/tony-tvu/goexpense/app"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/models"
	"github.com/tony-tvu/goexpense/util"
	"go.mongodb.org/mongo-driver/bson"
)

type UserHandler struct {
	App *app.App
}

var encryptionKey string
var jwtKey string

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("no .env file found")
	}
	encryptionKey = os.Getenv("ENCRYPTION_KEY")
	jwtKey = os.Getenv("JWT_KEY")
	if util.ContainsEmpty(encryptionKey, jwtKey) {
		log.Fatal("auth keys are missing")
	}
}

func (h UserHandler) GetInfo(c *gin.Context) {
	ctx := c.Request.Context()
	cookie, err := c.Request.Cookie("goexpense_access")

	// no access token - make user log in
	if err != nil || cookie.Value == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	type Claims struct {
		Email string
		jwt.RegisteredClaims
	}

	decrypted, err := auth.Decrypt(encryptionKey, cookie.Value)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	/*
		No need to check if access token is valid/expired because LoggedIn middleware
		has already validated it and might've refreshed the access token already, which
		would be in w http.ResponseWriter. We only care about the email from the
		original r *http.Request here.
	*/
	tkn, _ := jwt.ParseWithClaims(decrypted, &Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})
	claims := tkn.Claims.(*Claims)

	if claims.Email == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var u *models.User
	err = h.App.Users.FindOne(ctx, bson.D{{Key: "email", Value: claims.Email}}).Decode(&u)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	// do not send back hashed password
	u.Password = ""
	c.JSON(200, u)
}

// Send email invite to new user
func (h UserHandler) Invite(c *gin.Context) {
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
	log.Println("GetSessions called")
}
