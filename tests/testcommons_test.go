package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tony-tvu/goexpense/models"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type QLResponse struct {
	Data   json.RawMessage
	Errors []struct{ Message string }
}

// Log user in and return tokens
func logUserIn(t *testing.T, username, password string) (string, string, int) {
	t.Helper()
	m := map[string]string{"username": username, "password": password}
	b := new(bytes.Buffer)

	err := json.NewEncoder(b).Encode(m)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/login", srv.URL), b)
	require.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	if res.StatusCode != 200 {
		return "", "", res.StatusCode
	}

	cookies := getCookies(t, res.Cookies())
	accessToken := cookies["goexpense_access"]
	if accessToken == "" {
		t.FailNow()
	}
	refreshToken := cookies["goexpense_refresh"]
	if refreshToken == "" {
		t.FailNow()
	}
	return accessToken, refreshToken, res.StatusCode
}

// Save a new user to db
func createTestUser(t *testing.T) (*models.User, func()) {
	t.Helper()

	username := fmt.Sprint(time.Now().UnixNano())
	password := "password123!"

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)

	user := &models.User{
		Username: username,
		Email:    fmt.Sprintf("%v@email.com", username),
		Password: string(hash),
		Type:     models.RegularUser,
	}

	doc := &bson.D{
		{Key: "username", Value: username},
		{Key: "email", Value: fmt.Sprintf("%v@email.com", username)},
		{Key: "password", Value: string(hash)},
		{Key: "type", Value: models.RegularUser},
		{Key: "created_at", Value: time.Now()},
		{Key: "updated_at", Value: time.Now()},
	}
	if _, err = testApp.Db.Users.InsertOne(ctx, doc); err != nil {
		t.FailNow()
	}

	user.Password = password
	return user, func() {
		deleteUser(t, username)
	}
}

func deleteUser(t *testing.T, username string) {
	t.Helper()


	_, err := testApp.Db.Users.DeleteOne(ctx, bson.M{"username": username})
	if err != nil {
		t.FailNow()
	}
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
func makeRequest(t *testing.T, method string, url string, accessToken *string, refreshToken *string, body ...map[string]string) (res *http.Response) {
	t.Helper()

	var req *http.Request

	if len(body) > 0 {
		b := body[0]
		bodyJSON, err := json.Marshal(b)
		require.NoError(t, err)
		req, _ = http.NewRequest(method, fmt.Sprintf("%s%s", srv.URL, url), bytes.NewBuffer(bodyJSON))
	} else {
		req, _ = http.NewRequest(method, fmt.Sprintf("%s%s", srv.URL, url), nil)
	}

	if accessToken != nil {
		req.AddCookie(&http.Cookie{
			Name:  "goexpense_access",
			Value: *accessToken})
	}

	if refreshToken != nil {
		req.AddCookie(&http.Cookie{
			Name:  "goexpense_refresh",
			Value: *refreshToken})
	}

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return res
}
