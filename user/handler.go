package user

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/tony-tvu/goexpense/app"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	App *app.App
}

func (h UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var u models.User
	err := json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Save new user
	doc := &bson.D{
		{Key: "email", Value: u.Email},
		{Key: "name", Value: u.Name},
		{Key: "password", Value: string(hash)},
		{Key: "role", Value: models.ExternalUser},
		{Key: "verified", Value: false},
		{Key: "created_at", Value: time.Now()},
	}
	_, err = h.App.Users.InsertOne(ctx, doc)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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
		UserID string `json:"user_id"`
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
		would be in w http.ResponseWriter. We only care about the user_id from the
		original r *http.Request here.
	*/
	tkn, _ := jwt.ParseWithClaims(decrypted, &Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(h.App.JwtKey), nil
		})
	claims := tkn.Claims.(*Claims)

	var u *models.User
	objID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.App.Users.FindOne(ctx, bson.D{{Key: "_id", Value: objID}}).Decode(&u)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// do not send back hashed password
	u.Password = ""
	json.NewEncoder(w).Encode(u)

}
