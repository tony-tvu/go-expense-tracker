package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/tony-tvu/goexpense/app"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/models"
	"go.mongodb.org/mongo-driver/bson"
)

type Input struct {
	Email    string
	Password string
}

func LoginEmail(ctx context.Context, a *app.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var input Input
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Find existing user
		var user *models.User
		coll := a.MongoClient.Database(a.Db).Collection(a.Coll.Users)
		err = coll.FindOne(ctx, bson.D{{Key: "email", Value: input.Email}}).Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Verify password
		decrypted, _ := auth.Decrypt(a.Secret, user.Password)
		if input.Password != decrypted {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Println("success")
		// create refresh token
		refreshToken, accessToken, err := auth.GetAuthTokens(a, user)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Println(refreshToken)
		log.Println(accessToken)
		// create access token

		// return refresh and access tokens
	}
}

func IsTokenValid(ctx context.Context, a *app.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		type Tokens struct {
			RefreshToken string
			AccessToken  string
		}

		var tokens Tokens
		err := json.NewDecoder(r.Body).Decode(&tokens)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		isValid := auth.IsTokenValid(a, tokens.RefreshToken)

		log.Println(isValid)
	}
}
