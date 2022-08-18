package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/tony-tvu/goexpense/app"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/models"

	"go.mongodb.org/mongo-driver/bson"
)

func UserHandler(a *app.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			err := createUser(a, w, r)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			fmt.Fprint(w, "User created")
			return
		} else {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
	}
}

func createUser(a *app.App, w http.ResponseWriter, r *http.Request) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var u models.User
	err := json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		return errors.New("decode error")
	}

	// Encrypt password
	encrypted, err := auth.Encrypt(a.AuthKey, u.Password)
	if err != nil {
		return errors.New("encrypt error")
	}

	// Save new user
	coll := a.MongoClient.Database(a.DbName).Collection(a.UserCollection)
	doc := bson.D{
		{Key: "email", Value: u.Email},
		{Key: "name", Value: u.Name},
		{Key: "password", Value: encrypted},
		{Key: "role", Value: models.ExternalUser},
		{Key: "verified", Value: false},
	}

	_, err = coll.InsertOne(ctx, doc)
	if err != nil {
		return errors.New("db request error")
	}

	return nil
}
