package tests

import (
	"context"
	"log"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/tony-tvu/goexpense/app"
	"github.com/tony-tvu/goexpense/server"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	s               *server.Server
	srv             *httptest.Server
	ctx             context.Context
	refreshTokenExp int
	accessTokenExp  int
)

func TestMain(m *testing.M) {
	// BeforeAll
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()

	refreshTokenExp = 2
	accessTokenExp = 1

	s = &server.Server{
		App: &app.App{
			Env:               "test",
			Port:              "5000",
			EncryptionKey:     "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			JwtKey:            "jwt_key",
			RefreshTokenExp:   refreshTokenExp,
			AccessTokenExp:    accessTokenExp,
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

	mongoclient, err := mongo.Connect(ctx, options.Client().ApplyURI(s.App.MongoURI))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := mongoclient.Disconnect(ctx); err != nil {
			log.Println("mongo has been disconnected: ", err)
		}
	}()

	s.App.Users = mongoclient.Database(s.App.DbName).Collection("users")
	s.App.Sessions = mongoclient.Database(s.App.DbName).Collection("sessions")
	s.App.Users.Drop(ctx)
	s.App.Sessions.Drop(ctx)

	srv = httptest.NewServer(s.App.Router)

	// Run tests
	exitVal := m.Run()

	// Teardown
	os.Exit(exitVal)
}
