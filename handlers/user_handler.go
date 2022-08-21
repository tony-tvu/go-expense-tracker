package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/tony-tvu/goexpense/app"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/models"
	"go.mongodb.org/mongo-driver/bson"
)

type Credentials struct {
	Email    string
	Password string
}

func CreateUser(ctx context.Context, a *app.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var u models.User
		err := json.NewDecoder(r.Body).Decode(&u)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Encrypt password
		encrypted, err := auth.Encrypt(a.Secret, u.Password)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Save new user
		coll := a.MongoClient.Database(a.Db).Collection(a.Coll.Users)
		doc := bson.D{
			{Key: "email", Value: u.Email},
			{Key: "name", Value: u.Name},
			{Key: "password", Value: encrypted},
			{Key: "role", Value: models.ExternalUser},
			{Key: "verified", Value: false},
			{Key: "created_at", Value: time.Now()},
		}

		_, err = coll.InsertOne(ctx, doc)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, "User created")
	}
}

func LoginEmail(ctx context.Context, a *app.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var c Credentials
		err := json.NewDecoder(r.Body).Decode(&c)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// find existing user
		var u *models.User
		coll := a.MongoClient.Database(a.Db).Collection(a.Coll.Users)
		err = coll.FindOne(ctx, bson.D{{Key: "email", Value: c.Email}}).Decode(&u)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// verify password
		decrypted, _ := auth.Decrypt(a.Secret, u.Password)
		if c.Password != decrypted {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// create refresh token
		refreshToken, err := auth.CreateRefreshToken(ctx, a, u)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		coll = a.MongoClient.Database(a.Db).Collection(a.Coll.Sessions)

		// delete existing sessions
		_, err = coll.DeleteOne(ctx, bson.M{"user_id": u.ObjectID.Hex()})
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

		_, err = coll.InsertOne(ctx, doc)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// create access token
		accessToken, err := auth.CreateAccessToken(ctx, a, u.ObjectID.Hex(), string(u.Role))
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

func GetUserInfo(ctx context.Context, a *app.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "User info")

	}
}
