package smoketest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tony-tvu/goexpense/app"
	"github.com/tony-tvu/goexpense/models"
	"github.com/tony-tvu/goexpense/server"
	"github.com/tony-tvu/goexpense/user"
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

// Log user in and return access token
func logUserIn(t *testing.T, email, password string) (string, int) {
	t.Helper()
	m, b := map[string]string{
		"email":    email,
		"password": password},
		new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(m)
	require.NoError(t, err)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/login", srv.URL), b)
	require.NoError(t, err)

	res, _ := client.Do(req)
	if res.StatusCode != 200 {
		return "", res.StatusCode
	}

	cookies := getCookies(t, res.Cookies())
	access_token := cookies["goexpense_access"]
	if access_token == "" {
		t.FailNow()
	}
	return access_token, res.StatusCode
}

// Save a new user to db
func createUser(t *testing.T, a *app.App, name, email, password string) {
	err := user.SaveUser(context.TODO(), a, &models.User{
		Name:     name,
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
}

// Return cookies map from http response cookies
func getCookies(t *testing.T, cookies_res []*http.Cookie) map[string]string {
	t.Helper()
	cookies := make(map[string]string)
	for _, cookie := range cookies_res {
		cookies[cookie.Name] = cookie.Value
	}
	return cookies
}
