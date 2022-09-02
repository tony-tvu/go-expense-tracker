package smoketest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tony-tvu/goexpense/database"
	"github.com/tony-tvu/goexpense/models"
	"github.com/tony-tvu/goexpense/user"
)

// Log user in and return access token
func logUserIn(t *testing.T, email, password string) (string, int) {
	t.Helper()
	m := map[string]string{"email": email, "password": password}
	b := new(bytes.Buffer)

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
	accessToken := cookies["goexpense_access"]
	if accessToken == "" {
		t.FailNow()
	}
	return accessToken, res.StatusCode
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

func MakeApiRequest(t *testing.T, method string, url string, body *map[string]string, accessToken *string) (res *http.Response) {
	t.Helper()
	client := &http.Client{}
	var req *http.Request

	var jsonBytes *bytes.Buffer
	if body != nil {
		jsonBytes = new(bytes.Buffer)
		err := json.NewEncoder(jsonBytes).Encode(body)
		require.NoError(t, err)

		req, _ = http.NewRequest(method, fmt.Sprintf("%s%s", srv.URL, url), jsonBytes)
	} else {
		req, _ = http.NewRequest(method, fmt.Sprintf("%s%s", srv.URL, url), nil)
	}

	if accessToken != nil {
		req.AddCookie(&http.Cookie{
			Name:  "goexpense_access",
			Value: *accessToken})
	}

	res, err := client.Do(req)
	require.NoError(t, err)
	return res
}
