package tests

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tony-tvu/goexpense/entity"
)

// VerifyUsers should issue access tokens correctly
func TestAuthTokens(t *testing.T) {
	t.Parallel()

	// create user and login
	username := fmt.Sprint(time.Now().UnixNano())
	email := fmt.Sprintf("%v@email.com", username)
	password := "password"
	_, cleanup := createUser(t, username, email, password)
	defer cleanup()
	accessToken, refreshToken, _ := logUserIn(t, username, password)

	// make request to endpoint where user must be logged in
	res := makeRequest(t, "GET", "/api/user_info", &accessToken, &refreshToken)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// should not return a new access token
	cookies := getCookies(t, res.Cookies())
	assert.Equal(t, "", cookies["goexpense_access"])

	// make request with mocked expired access token
	res = makeRequest(t, "GET", "/api/user_info", nil, &refreshToken)

	// should return new access token and have no errors
	assert.Equal(t, http.StatusOK, res.StatusCode)
	cookies = getCookies(t, res.Cookies())
	assert.NotEqual(t, accessToken, cookies["goexpense_access"])
}

// VerifyUsers should return 401 for revoked sessions
func TestAuthRevokeTokens(t *testing.T) {
	t.Parallel()

	// create user and login
	username := fmt.Sprint(time.Now().UnixNano())
	email := fmt.Sprintf("%v@email.com", username)
	password := "password"
	user, cleanup := createUser(t, username, email, password)
	defer cleanup()
	accessToken, refreshToken, _ := logUserIn(t, username, password)

	// make request to endpoint where user must be logged in
	res := makeRequest(t, "GET", "/api/user_info", &accessToken, &refreshToken)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// make same request with expired accessToken
	res = makeRequest(t, "GET", "/api/user_info", nil, &refreshToken)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// revoke token
	if result := testApp.Db.Exec("DELETE FROM sessions WHERE user_id = ?", user.ID); result.Error != nil {
		t.FailNow()
	}

	// make request to authRequired endpoint with expired accessToken
	res = makeRequest(t, "GET", "/api/user_info", nil, &refreshToken)

	// 	// should return 401
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

// Auth required endpoints should not allow access to guest users
func TestAuthRestrictAccess(t *testing.T) {
	t.Parallel()

	// make request to endpoint where user must be logged in without access token
	res := makeRequest(t, "GET", "/api/user_info", nil, nil)

	// should return 401
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

	// response should not have new token cookie
	cookies := getCookies(t, res.Cookies())
	assert.Equal(t, "", cookies["goexpense_access"])
	assert.Equal(t, "", cookies["goexpense_refresh"])

	// make request to same endpoint with invalid access token
	invalidToken := "invalidToken"
	res = makeRequest(t, "GET", "/api/user_info", &invalidToken, nil)

	// should return 401
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

	// response should not have new token cookie
	cookies = getCookies(t, res.Cookies())
	assert.Equal(t, "", cookies["goexpense_access"])
	assert.Equal(t, "", cookies["goexpense_refresh"])
}

// Admin restricted handlers should not allow guest or regular users
func TestAuthAdminRestricted(t *testing.T) {
	t.Parallel()

	// make request to admin-only route as a guest user
	res := makeRequest(t, "GET", "/api/sessions", nil, nil)

	// should return 401
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

	// create regular user and login
	username := fmt.Sprint(time.Now().UnixNano())
	email := fmt.Sprintf("%v@email.com", username)
	password := "password"
	_, cleanup := createUser(t, username, email, password)
	defer cleanup()
	accessToken, refreshToken, _ := logUserIn(t, username, password)

	// make request to admin-only route as regular user
	res = makeRequest(t, "GET", "/api/sessions", &refreshToken, &accessToken)

	// should return 401
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

	// logout user
	res = makeRequest(t, "POST", "/api/logout", &accessToken, &refreshToken)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// make user an admin and login
	if result := testApp.Db.Model(&entity.User{}).Where("username = ?", username).Update("type", entity.AdminUser); result.Error != nil {
		t.FailNow()
	}

	accessToken, refreshToken, _ = logUserIn(t, username, password)

	// make request to admin-only endpoint as admin
	res = makeRequest(t, "GET", "/api/sessions", &accessToken, &refreshToken)

	// should return 200
	assert.Equal(t, http.StatusOK, res.StatusCode)
}
