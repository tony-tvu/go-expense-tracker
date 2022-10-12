package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tony-tvu/goexpense/types"
	"go.mongodb.org/mongo-driver/bson"
)

// AuthorizeUser should issue access tokens correctly
func TestAuthTokens(t *testing.T) {
	t.Parallel()

	// create user and login
	testUser, cleanup := createTestUser(t)
	defer cleanup()
	accessToken, refreshToken, _ := logUserIn(t, testUser.Username, testUser.Password)

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

// AuthorizeUser should return 401 for revoked sessions
func TestAuthRevokeTokens(t *testing.T) {
	t.Parallel()

	// create user and login
	testUser, cleanup := createTestUser(t)
	defer cleanup()
	accessToken, refreshToken, _ := logUserIn(t, testUser.Username, testUser.Password)

	// make request to endpoint where user must be logged in
	res := makeRequest(t, "GET", "/api/user_info", &accessToken, &refreshToken)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// make same request with expired accessToken
	res = makeRequest(t, "GET", "/api/user_info", nil, &refreshToken)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// revoke token
	_, err := testApp.Db.Sessions.DeleteMany(ctx, bson.M{"user_id": testUser.ID})
	if err != nil {
		t.FailNow()
	}

	// make request to authRequired endpoint with expired accessToken
	res = makeRequest(t, "GET", "/api/user_info", nil, &refreshToken)

	// should return 401
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
	user, cleanup := createTestUser(t)
	defer cleanup()
	accessToken, refreshToken, _ := logUserIn(t, user.Username, user.Password)

	// make request to admin-only route as regular user
	res = makeRequest(t, "GET", "/api/sessions", &refreshToken, &accessToken)

	// should return 401
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

	// logout user
	res = makeRequest(t, "POST", "/api/logout", &accessToken, &refreshToken)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// make user an admin and login
	_, err := testApp.Db.Users.UpdateOne(
		ctx,
		bson.M{"username": user.Username},
		bson.M{
			"$set": bson.M{
				"user_type": types.AdminUser,
			}},
	)
	if err != nil {
		t.FailNow()
	}

	accessToken, refreshToken, _ = logUserIn(t, user.Username, user.Password)

	// make request to admin-only endpoint as admin
	res = makeRequest(t, "GET", "/api/sessions", &accessToken, &refreshToken)

	// should return 200
	assert.Equal(t, http.StatusOK, res.StatusCode)
}
