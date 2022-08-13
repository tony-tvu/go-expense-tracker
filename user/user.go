package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/tony-tvu/goexpense/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserConfigs struct {
	Client     *mongo.Client
	Database   string
	Collection string
	AuthKey    []byte
	DBTimeout  int
}

type Role string

const (
	Admin    Role = "Admin"
	External Role = "External"
)

type User struct {
	Name     string
	Email    string
	Password string
	Role     Role
	Verified bool
}

func (conf *UserConfigs) Handler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.DBTimeout)*time.Second)
	defer cancel()
	
	if r.Method == "POST" {
		createUser(ctx, conf, w, r)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func createUser(ctx context.Context, conf *UserConfigs, w http.ResponseWriter, r *http.Request) {
	var u User
	err := json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		w.WriteHeader(400)
		fmt.Fprint(w, "Bad Request")
		return
	}

	// Encrypt password
	encrypted, err := auth.Encrypt(conf.AuthKey, u.Password)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "Bad Request")
		return
	}

	// Save new user
	coll := conf.Client.Database(conf.Database).Collection(conf.Collection)
	doc := bson.D{
		{Key: "email", Value: u.Email},
		{Key: "name", Value: u.Name},
		{Key: "password", Value: encrypted},
		{Key: "role", Value: External},
		{Key: "verified", Value: false},
	}

	_, err = coll.InsertOne(ctx, doc)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "Bad Request")
		return
	}

	w.WriteHeader(200)
	fmt.Fprint(w, "User created")
}
