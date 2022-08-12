package user

import (
	"context"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserConfigs struct {
	Client     *mongo.Client
	Database   string
	Collection string
}

type User struct {
	Name     string
	Email    string
	Password string
}

func (s *UserConfigs) Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		createUser(s, w, r)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func createUser(s *UserConfigs, w http.ResponseWriter, r *http.Request) {
	var u User

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	coll := s.Client.Database(s.Database).Collection(s.Collection)
	doc := bson.D{
		{Key: "email", Value: u.Email},
		{Key: "name", Value: u.Name},
		{Key: "password", Value: u.Password},
	}

	_, err = coll.InsertOne(context.TODO(), doc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

}
