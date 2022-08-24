package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tony-tvu/goexpense/app"
	"github.com/tony-tvu/goexpense/server"
	"github.com/tony-tvu/goexpense/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	s     *server.Server
	name  string
	email string
)

func TestMain(m *testing.M) {
	// BeforeAll
	s = &server.Server{
		App: &app.App{
			Env:               "test",
			Port:              "5000",
			JwtKey:            "jwt_key",
			RefreshTokenExp:   10,
			AccessTokenExp:    5,
			MongoURI:          "mongodb://localhost:27017/local_db",
			DbName:            "goexpense_test",
			PlaidClientID:     "plaidClientID",
			PlaidSecret:       "plaidSecret",
			PlaidEnv:          "sandbox",
			PlaidCountryCodes: "US,CA",
			PlaidProducts:     "auth,transactions",
		},
	}
	s.Initialize()

	mongoclient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(s.App.MongoURI))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := mongoclient.Disconnect(context.TODO()); err != nil {
			log.Println("mongo has been disconnected: ", err)
		}
	}()
	s.App.Collections = &app.Collections{
		Users:    mongoclient.Database(s.App.DbName).Collection("users"),
		Sessions: mongoclient.Database(s.App.DbName).Collection("sessions"),
	}

	name = "Test"
	email = "test@email.com"

	// Run tests
	exitVal := m.Run()

	// Teardown
	os.Exit(exitVal)
}

// User needs to be logged in to access GetUserInfo
func TestGetUserInfo(t *testing.T) {
	// setup
	s.App.Collections.Users.Drop(context.TODO())
	s.App.Collections.Sessions.Drop(context.TODO())

	s.App.Collections.Users.InsertOne(context.TODO(), bson.D{
		{Key: "email", Value: email},
		{Key: "name", Value: name},
		{Key: "password",
			// Encrypted password
			Value: "66777bfc0f53772bb2f97e4bb5a4746d80f7bb4a0f416aade824a5001da43452393faf5b"},
		{Key: "role", Value: user.ExternalUser},
		{Key: "verified", Value: true},
		{Key: "created_at", Value: time.Now()},
	})
	writer := httptest.NewRecorder()

	// log user in
	m, body := map[string]string{
		"email":    "test@email.com",
		"password": "password"},
		new(bytes.Buffer)
	json.NewEncoder(body).Encode(m)

	loginHandler := http.HandlerFunc(s.App.Handlers.LoginEmail)
	loginReq := httptest.NewRequest(http.MethodPost, "/", body)
	loginHandler.ServeHTTP(writer, loginReq)

	assert.Equal(t, http.StatusOK, writer.Code)

	// infoHandler := http.HandlerFunc(s.App.Handlers.GetUserInfo)
	// infoReq := httptest.NewRequest(http.MethodGet, "/", nil)

	// // when
	// infoHandler.ServeHTTP(writer, infoReq)

	// // then: status is OK
	// assert.Equal(t, http.StatusOK, writer.Code)

	// // and: body has correct data
	// type Body struct {
	// 	Message string `json:"message"`
	// }

	// var b Body
	// json.NewDecoder(writer.Body).Decode(&b)
	// assert.Equal(t, "test@email.com", b.Message)
}

// func TestEmailLogin(t *testing.T) {
// 	// setup
// 	s.App.Collections.Users.Drop(context.TODO())
// 	s.App.Collections.Sessions.Drop(context.TODO())
// 	writer := httptest.NewRecorder()


// }
