package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/tony-tvu/goexpense/app"
	"github.com/tony-tvu/goexpense/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/api/idtoken"
)

type Session struct {
	ObjectID     primitive.ObjectID `bson:"_id" json:"_id"`
	UserID       string             `bson:"user_id"`
	RefreshToken string             `bson:"refresh_token"`
	CreatedAt    time.Time          `bson:"created_at"`
	ExpiresAt    time.Time          `bson:"expires_at"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func GoogleLogin(a *app.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		idToken := r.Header.Get("Google-ID-Token")

		if idToken == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		payload, err := idtoken.Validate(ctx, idToken, a.GoogleClientID)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		fmt.Print(payload.Claims)
	}
}

func EmailLogin(a *app.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var c Credentials
		err := json.NewDecoder(r.Body).Decode(&c)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// find existing user
		var u *user.User
		err = a.Collections.Users.FindOne(ctx, bson.D{{Key: "email", Value: c.Email}}).Decode(&u)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// verify password
		decrypted, _ := Decrypt(a.EncryptKey, u.Password)
		if c.Password != decrypted {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// create refresh token
		refreshToken, err := CreateRefreshToken(ctx, a, u)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// delete existing sessions
		_, err = a.Collections.Sessions.DeleteOne(ctx, bson.M{"user_id": u.ObjectID.Hex()})
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

		_, err = a.Collections.Sessions.InsertOne(ctx, doc)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// create access token
		accessToken, err := CreateAccessToken(ctx, a, u.ObjectID.Hex(), string(u.Role))
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
}
