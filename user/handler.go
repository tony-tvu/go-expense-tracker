package user

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/tony-tvu/goexpense/app"
	"go.mongodb.org/mongo-driver/bson"
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

		// Encrypt password
		// encrypted, err := auth.Encrypt(a.Secret, u.Password)
		// if err != nil {
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	return
		// }

		// Save new user
		doc := &bson.D{
			{Key: "email", Value: u.Email},
			{Key: "name", Value: u.Name},
			{Key: "role", Value: ExternalUser},
			{Key: "verified", Value: false},
			{Key: "created_at", Value: time.Now()},
		}
		_, err = a.Collections.Users.InsertOne(ctx, doc)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func GetInfo(a *app.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var u *User
		err := a.Collections.Users.FindOne(ctx, bson.D{{Key: "email", Value: "test@email.com"}}).Decode(&u)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		body := make(map[string]string)
		body["message"] = u.Email
		jData, _ := json.Marshal(body)
		w.Write(jData)
	}
}
