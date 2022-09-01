package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"github.com/tony-tvu/goexpense/app"
	"github.com/tony-tvu/goexpense/database"
	"github.com/tony-tvu/goexpense/models"
	"github.com/tony-tvu/goexpense/user"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	testApp         *app.App
	srv             *httptest.Server
	ctx             context.Context
	refreshTokenExp int
	accessTokenExp  int
)

func TestMain(m *testing.M) {
	ctx = context.Background()

	if err := godotenv.Load(".env"); err != nil {
		log.Println("no .env file found")
	}

	refreshExp, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXP"))
	if err != nil {
		log.Fatal(err)
	}
	refreshTokenExp = refreshExp

	accessExp, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXP"))
	if err != nil {
		log.Fatal(err)
	}
	accessTokenExp = accessExp

	mongoclient, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := mongoclient.Disconnect(ctx); err != nil {
			log.Println("mongo has been disconnected: ", err)
		}
	}()

	testApp = &app.App{}
	testApp.Initialize(ctx)
	testApp.Db.Users = mongoclient.Database("goexpense_test").Collection("users")
	testApp.Db.Sessions = mongoclient.Database("goexpense_test").Collection("sessions")

	// drop all collections
	testApp.Db.Users.Drop(ctx)
	testApp.Db.Sessions.Drop(ctx)

	// start test server
	srv = httptest.NewServer(testApp.Router)

	// run tests
	exitVal := m.Run()

	// teardown
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
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/login", srv.URL), b)
	require.NoError(t, err)

	res, err := client.Do(req)
	require.NoError(t, err)
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
func createUser(t *testing.T, db *database.Db, name, email, password string) {
	t.Helper()
	err := user.SaveUser(ctx, db, &models.User{
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
