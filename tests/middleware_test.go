package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tony-tvu/goexpense/entity"
)

func TestMiddlware(t *testing.T) {
	t.Run("AuthRequired middleware should issue access tokens correctly", func(t *testing.T) {
		t.Parallel()

		// create user and login
		username := "middleware"
		email := "middleware@email.com"
		password := "password"
		cleanup := createUser(t, username, email, password)
		defer cleanup()
		accessToken, refreshToken, _ := logUserIn(t, username, password)

		// make request to authRequired endpoint
		res := MakeApiRequest(t, "GET", "/api/user_info", &accessToken, &refreshToken)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		// should not return a new access token
		cookies := getCookies(t, res.Cookies())
		assert.Equal(t, "", cookies["goexpense_access"])

		// make request with mocked expired access token
		res = MakeApiRequest(t, "GET", "/api/user_info", nil, &refreshToken)

		// should return 200 with new access token
		assert.Equal(t, http.StatusOK, res.StatusCode)
		cookies = getCookies(t, res.Cookies())
		assert.NotEqual(t, accessToken, cookies["goexpense_access"])
	})

	t.Run("AuthRequired middleware should return 401 for revoked sessions", func(t *testing.T) {
		t.Parallel()

		// create user and login
		username := "revokeMe"
		email := "revokeMe@email.com"
		password := "^%#(GY%H=G$%asdf"
		cleanup := createUser(t, username, email, password)
		defer cleanup()
		accessToken, refreshToken, _ := logUserIn(t, username, password)

		// make request to authRequired endpoint
		res := MakeApiRequest(t, "GET", "/api/user_info", &accessToken, &refreshToken)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		// make request to authRequired endpoint with expired accessToken
		res = MakeApiRequest(t, "GET", "/api/user_info", nil, &refreshToken)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		// revoke token
		if result := testApp.Db.Exec("DELETE FROM sessions WHERE username = ?", username); result.Error != nil {
			t.FailNow()
		}

		// make request to authRequired endpoint with expired accessToken
		res = MakeApiRequest(t, "GET", "/api/user_info", nil, &refreshToken)

		// should return 401
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	})

	t.Run("AuthRequired middleware should not allow access to guest users", func(t *testing.T) {
		t.Parallel()

		// make request to endpoint where user must be logged in without access token
		res := MakeApiRequest(t, "GET", "/api/user_info", nil, nil)

		// should return 401
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

		// make request to same endpoint with invalid access token
		invalidToken := "invalidToken"
		res = MakeApiRequest(t, "GET", "/api/user_info", &invalidToken, nil)

		// should return 401
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

		// response should not have new token cookie
		cookies := getCookies(t, res.Cookies())
		assert.Equal(t, "", cookies["goexpense_access"])
		assert.Equal(t, "", cookies["goexpense_refresh"])
	})

	t.Run("AdminRequired middleware should not allow access to guest or regular users", func(t *testing.T) {
		t.Parallel()

		// make request to admin-only route as a guest user
		res := MakeApiRequest(t, "GET", "/api/sessions", nil, nil)

		// should return 401
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

		// create regular user and login
		username := "middleware3"
		email := "middleware3@email.com"
		password := "^%#(GY%H=G$%asdf"
		cleanup := createUser(t, username, email, password)
		defer cleanup()
		accessToken, refreshToken, _ := logUserIn(t, username, password)

		// make request to admin-only route as regular user
		res = MakeApiRequest(t, "GET", "/api/sessions", &refreshToken, &accessToken)

		// should return 401
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

		// logout user
		res = MakeApiRequest(t, "POST", "/api/logout", &accessToken, &refreshToken)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		// make user an admin and login
		if result := testApp.Db.Model(&entity.User{}).Where("username = ?", username).Update("type", entity.AdminUser); result.Error != nil {
			t.FailNow()
		}

		accessToken, refreshToken, _ = logUserIn(t, username, password)

		// make request to admin-only endpoint as admin
		res = MakeApiRequest(t, "GET", "/api/sessions", &accessToken, &refreshToken)

		// should return 200
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}
