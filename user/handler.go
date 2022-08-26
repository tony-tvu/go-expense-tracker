package user

import (
	"encoding/json"
	"log"
	"net/http"
	"net/mail"
	"net/smtp"
	"os"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"github.com/tony-tvu/goexpense/app"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/models"
	"go.mongodb.org/mongo-driver/bson"
)

type UserHandler struct {
	App *app.App
}

func (h UserHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cookie, err := r.Cookie("goexpense_access")

	// no access token - make user log in
	if err != nil || cookie.Value == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	type Claims struct {
		Email string
		jwt.RegisteredClaims
	}

	decrypted, err := auth.Decrypt(h.App.EncryptionKey, cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
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
			return []byte(h.App.JwtKey), nil
		})
	claims := tkn.Claims.(*Claims)

	if claims.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var u *models.User
	err = h.App.Users.FindOne(ctx, bson.D{{Key: "email", Value: claims.Email}}).Decode(&u)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// do not send back hashed password
	u.Password = ""
	json.NewEncoder(w).Encode(u)
}

// Send email invite to new user
func (h UserHandler) Invite(w http.ResponseWriter, r *http.Request) {
	type Body struct {
		Email string `json:"email"`
	}

	var b Body
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// validate email address
	_, err = mail.ParseAddress(b.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	to := []string{b.Email}

	godotenv.Load(".env")
	from := os.Getenv("EMAIL_SENDER")
	password := os.Getenv("EMAIL_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	if from == "" || password == "" || smtpHost == "" || smtpPort == "" {
		log.Println("error - email sender configs are missing")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	auth := smtp.PlainAuth("", from, password, smtpHost)

	subject := "Subject: This is the subject of the mail\n"
	body := "This is the body of the mail"
	message := []byte(subject + body)

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
