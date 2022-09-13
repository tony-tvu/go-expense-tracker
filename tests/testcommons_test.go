package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tony-tvu/goexpense/graph/models"
	"golang.org/x/crypto/bcrypt"
)

type QLResponse struct {
	Data   json.RawMessage
	Errors []struct{ Message string }
}

// Log user in and return access token
func logUserIn(t *testing.T, username, password string) (string, string, *QLResponse) {
	t.Helper()

	query := fmt.Sprintf(
		`mutation {
			login(
				input: {
					username: "%s"
					password: "%s"
				}
			)
		}`, username, password,
	)

	query = strings.Replace(query, "\t", "", -1)
	q := struct{ Query string }{Query: query}

	data, err := json.Marshal(q)
	if err != nil {
		t.Fatal("failed to marshal graphql query")
	}

	url := srv.URL + "/api/graphql"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, 200, res.StatusCode)

	var qlRes QLResponse
	err = json.NewDecoder(res.Body).Decode(&qlRes)
	if err != nil {
		t.Fatal("failed to parse GraphQL response:", err)
	}
	if len(qlRes.Errors) > 0 {
		return "", "", &qlRes
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
	return accessToken, refreshToken, &qlRes
}

// Save a new user to db
func createUser(t *testing.T, username, email, password string) func() {
	t.Helper()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)

	if result := testApp.Db.Create(&models.User{
		Username: username,
		Email:    email,
		Password: string(hash),
		Type:     models.UserTypeRegular,
	}); result.Error != nil {
		t.FailNow()
	}

	return func() {
		deleteUser(t, username)
	}
}

func deleteUser(t *testing.T, username string) {
	if result := testApp.Db.Exec("DELETE FROM users WHERE username = ?", username); result.Error != nil {
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

func doQL(t *testing.T, accessToken *string, refreshToken *string, query ...string) (*http.Response, *QLResponse) {
	t.Helper()
	var req *http.Request
	url := srv.URL + "/api/graphql"

	if len(query) > 0 {
		rawQuery := query[0]
		rawQuery = strings.Replace(rawQuery, "\t", "", -1)
		q := struct{ Query string }{Query: rawQuery}

		data, err := json.Marshal(q)
		if err != nil {
			t.Fatal("failed to marshal graphql query")
		}

		req, _ = http.NewRequest("POST", url, bytes.NewBuffer(data))
	} else {
		req, _ = http.NewRequest("POST", url, nil)
	}
	req.Header.Set("Content-Type", "application/json")

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
	defer res.Body.Close()

	var qlRes QLResponse
	err = json.NewDecoder(res.Body).Decode(&qlRes)
	if err != nil {
		t.Fatal("failed to parse GraphQL response:", err)
	}

	return res, &qlRes
}
