package user

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/tony-tvu/goexpense/app"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func Create(a *app.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var u User
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
			{Key: "role", Value: ExternalUser},
			{Key: "verified", Value: false},
			{Key: "created_at", Value: time.Now()},
		}
		_, err = a.Users.InsertOne(ctx, doc)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func GetInfo(a *app.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cookie, err := r.Cookie("goexpense_access")

		// no access token - make user log in
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		type Claims struct {
			UserID string `json:"user_id"`
			jwt.RegisteredClaims
		}

		tkn, _ := jwt.ParseWithClaims(cookie.Value, &Claims{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(a.JwtKey), nil
			})
		claims := tkn.Claims.(*Claims)

		var u *User
		objID, err := primitive.ObjectIDFromHex(claims.UserID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = a.Users.FindOne(ctx, bson.D{{Key: "_id", Value: objID}}).Decode(&u)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		// do not send back hashed password
		u.Password = ""
		json.NewEncoder(w).Encode(u)
	}
}
