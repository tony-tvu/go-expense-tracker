package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tony-tvu/goexpense/app"
	"github.com/tony-tvu/goexpense/models"
	"github.com/tony-tvu/goexpense/user"
)

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
