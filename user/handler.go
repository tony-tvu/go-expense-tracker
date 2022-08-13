package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/config"

	"go.mongodb.org/mongo-driver/bson"
)

type Handler struct {
	Config config.Config
}

func (h Handler) NewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := createUser(w, r, h.Config)
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

func createUser(w http.ResponseWriter, r *http.Request, cfg config.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.DBTimeout)*time.Second)
	defer cancel()

	var u User
	err := json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		return errors.New("decode error")
	}

	// Encrypt password
	encrypted, err := auth.Encrypt(cfg.AuthKey, u.Password)
	if err != nil {
		return errors.New("encrypt error")
	}

	// Save new user
	coll := cfg.Client.Database(cfg.Database).Collection(cfg.UserCollection)
	doc := bson.D{
		{Key: "email", Value: u.Email},
		{Key: "name", Value: u.Name},
		{Key: "password", Value: encrypted},
		{Key: "role", Value: External},
		{Key: "verified", Value: false},
	}

	_, err = coll.InsertOne(ctx, doc)
	if err != nil {
		return errors.New("db request error")
	}

	return nil
}
