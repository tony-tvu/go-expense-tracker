package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

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

func (h AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var c Credentials
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// find existing user
	var u *models.User
	err = h.App.Users.FindOne(ctx, bson.D{{Key: "email", Value: c.Email}}).Decode(&u)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// verify password
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(c.Password))
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// create refresh token
	refreshToken, err := CreateRefreshToken(ctx, h.App, u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// delete existing sessions
	_, err = h.App.Sessions.DeleteOne(ctx, bson.M{"user_id": u.ObjectID.Hex()})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// save new session
	doc := bson.D{
		{Key: "user_id", Value: u.ObjectID.Hex()},
		{Key: "refresh_token", Value: refreshToken.Value},
		{Key: "created_at", Value: time.Now()},
		{Key: "expires_at", Value: refreshToken.ExpiresAt},
	}

	_, err = h.App.Sessions.InsertOne(ctx, doc)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// create access token
	accessToken, err := CreateAccessToken(ctx, h.App, u.ObjectID.Hex(), string(u.Role))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "goexpense_access",
		Value:    accessToken.Value,
		Expires:  accessToken.ExpiresAt,
		HttpOnly: true,
	})

}

func (h AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cookie, err := r.Cookie("goexpense_access")

	// missing access_token
	if err != nil || cookie.Value == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// decrypt access token
	decrypted, err := Decrypt(h.App.EncryptionKey, cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tkn, err := jwt.ParseWithClaims(decrypted, &Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(h.App.JwtKey), nil
		})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	claims := tkn.Claims.(*Claims)

	_, err = h.App.Sessions.DeleteMany(ctx, bson.M{"user_id": claims.UserID})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (h AuthHandler) GetSessions(w http.ResponseWriter, r *http.Request) {
	log.Println("GetSessions called")
}
