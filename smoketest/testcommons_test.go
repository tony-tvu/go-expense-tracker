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
	"go.mongodb.org/mongo-driver/bson"
)

// Log user in and return access token
func logUserIn(t *testing.T, email, password string) (string, string, int) {
	t.Helper()
	m := map[string]string{"email": email, "password": password}
	b := new(bytes.Buffer)

	err := json.NewEncoder(b).Encode(m)
	require.NoError(t, err)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/login", srv.URL), b)
	require.NoError(t, err)

	res, err := client.Do(req)
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
func createUser(t *testing.T, db *database.Db, name, email, password string) func() {
	t.Helper()
	user, err := user.SaveUser(ctx, db, &models.User{
		Name:     name,
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)

	return func() {
		deleteUser(t, user)
	}
}

func deleteUser(t *testing.T, u *models.User) {
	_, err := testApp.Db.Users.DeleteOne(ctx, bson.M{"email": u.Email})
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

func MakeApiRequest(t *testing.T, method string, url string, accessToken *string, refreshToken *string, body ...map[string]string) (res *http.Response) {
	t.Helper()
	client := &http.Client{}
	var req *http.Request

	if len(body) > 0 {
		wtf := body[0]
		bodyJSON, err := json.Marshal(wtf)
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

	res, err := client.Do(req)
	require.NoError(t, err)
	return res
}
