package smoketest

import (
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
	"github.com/tony-tvu/goexpense/models"
	"github.com/tony-tvu/goexpense/server"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	s *server.Server
)

func TestMain(m *testing.M) {
	// Setup
	s = &server.Server{
		App: &app.App{
			Env:               "test",
			Port:              "5000",
			Secret:            "ThisKeyStringIs32BytesLongTest01",
			JwtKey:            "jwt_key",
			RefreshTokenExp:   10,
			AccessTokenExp:    5,
			MongoURI:          "mongodb://localhost:27017/local_db",
			Db:                "goexpense_test",
			PlaidClientID:     "plaidClientID",
			PlaidSecret:       "plaidSecret",
			PlaidEnv:          "sandbox",
			PlaidCountryCodes: "US,CA",
			PlaidProducts:     "auth,transactions",
		},
	}
	s.Initialize(context.TODO())

	mongoclient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(s.App.MongoURI))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := mongoclient.Disconnect(context.TODO()); err != nil {
			log.Println("mongo has been disconnected: ", err)
		}
	}()
	s.App.MongoClient = mongoclient

	// Run tests
	exitVal := m.Run()

	// Teardown
	os.Exit(exitVal)
}

func TestGetUserInfo(t *testing.T) {
	// given
	handler := http.HandlerFunc(s.App.Handlers.GetUserInfo)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	// req.Header.Add(server.HEADER_KEY_X_ACCOUNT, "myaccount")
	writer := httptest.NewRecorder()

	coll := s.App.MongoClient.Database(s.App.Db).Collection(s.App.Coll.Users)
	doc := bson.D{
		{Key: "email", Value: "test@email.com"},
		{Key: "name", Value: "test"},
		{Key: "password", Value: "password"},
		{Key: "role", Value: models.ExternalUser},
		{Key: "verified", Value: false},
		{Key: "created_at", Value: time.Now()},
	}

	_, err := coll.InsertOne(context.TODO(), doc)
	if err != nil {
		log.Fatal(err)
	}

	// when
	handler.ServeHTTP(writer, req)

	// then: status is OK
	assert.Equal(t, http.StatusOK, writer.Code)

	// and: body has correct data
	type Body struct {
		Message string `json:"message"`
	}

	var b Body
	err = json.NewDecoder(writer.Body).Decode(&b)
	if err != nil {
		log.Fatalln(err)
	}

	assert.Equal(t, "test@email.com", b.Message)
}
