package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tony-tvu/goexpense/entity"
)

// VerifyUsers should issue access tokens correctly
func TestAuthTokens(t *testing.T) {
	t.Parallel()

	// create user and login
	username := "isauthorized"
	email := "isauthorized@email.com"
	password := "password"
	_, cleanup := createUser(t, username, email, password)
	defer cleanup()
	accessToken, refreshToken, _ := logUserIn(t, username, password)

	// make request to authRequired endpoint
	query :=
		`query {
			userInfo {
				id
			}
		}`
	res, qlRes := doQL(t, &accessToken, &refreshToken, query)
	assert.Nil(t, qlRes.Errors)

	// should not return a new access token
	cookies := getCookies(t, res.Cookies())
	assert.Equal(t, "", cookies["goexpense_access"])

	// make request with mocked expired access token
	res, qlRes = doQL(t, nil, &refreshToken, query)

	// should return new access token and have no errors
	assert.Nil(t, qlRes.Errors)
	cookies = getCookies(t, res.Cookies())
	assert.NotEqual(t, accessToken, cookies["goexpense_access"])
}

// Resolvers should return error for revoked sessions
func TestAuthRevokeTokens(t *testing.T) {
	t.Parallel()

	// create user and login
	username := "revokeMe"
	email := "revokeMe@email.com"
	password := "^%#(GY%H=G$%asdf"
	user, cleanup := createUser(t, username, email, password)
	defer cleanup()
	accessToken, refreshToken, _ := logUserIn(t, username, password)

	// make request to authRequired endpoint
	query :=
		`query {
			userInfo {
				id
			}
		}`
	_, qlRes := doQL(t, &accessToken, &refreshToken, query)
	assert.Nil(t, qlRes.Errors)

	// make request to authRequired endpoint with expired accessToken
	_, qlRes = doQL(t, nil, &refreshToken, query)
	assert.Nil(t, qlRes.Errors)

	// revoke token
	if result := testApp.Db.Exec("DELETE FROM sessions WHERE user_id = ?", user.ID); result.Error != nil {
		t.FailNow()
	}

	// make request to authRequired endpoint with expired accessToken
	_, qlRes = doQL(t, nil, &refreshToken, query)

	// should return error
	assert.Equal(t, "not authorized", qlRes.Errors[0].Message)
}

// Auth restricted resolvers should not allow access to guest users
func TestAuthRestrictAccess(t *testing.T) {
	t.Parallel()

	// make request to endpoint where user must be logged in without access token
	query :=
		`query {
				userInfo {
					id
				}
			}`
	_, qlRes := doQL(t, nil, nil, query)

	// should return auth error
	assert.Equal(t, "not authorized", qlRes.Errors[0].Message)

	// make request to same endpoint with invalid access token
	invalidToken := "invalidToken"
	res, qlRes := doQL(t, &invalidToken, nil, query)

	// should return auth error
	assert.Equal(t, "not authorized", qlRes.Errors[0].Message)

	// response should not have new token cookie
	cookies := getCookies(t, res.Cookies())
	assert.Equal(t, "", cookies["goexpense_access"])
	assert.Equal(t, "", cookies["goexpense_refresh"])
}

// Admin restricted resolvers should not allow guest or regular users
func TestAuthAdminRestricted(t *testing.T) {
	t.Parallel()

	// make request to admin-only route as a guest user
	query :=
		`query {
			sessions {
				id
			}
		}`
	_, qlRes := doQL(t, nil, nil, query)

	// should return auth error
	assert.Equal(t, "not authorized", qlRes.Errors[0].Message)

	// create regular user and login
	username := "adminRestricted"
	email := "adminRestricted@email.com"
	password := "^%#(GY%H=G$%asdf"
	user, cleanup := createUser(t, username, email, password)
	defer cleanup()
	accessToken, refreshToken, _ := logUserIn(t, username, password)

	// make request to admin-only route as regular user
	_, qlRes = doQL(t, &accessToken, &refreshToken, query)

	// should return auth error
	assert.Equal(t, "not authorized", qlRes.Errors[0].Message)

	// logout user
	logoutQ :=
		`mutation {
			logout
		}`
	_, qlRes = doQL(t, &accessToken, &refreshToken, logoutQ)
	assert.Nil(t, qlRes.Errors)

	// make user an admin and login
	result := testApp.Db.Model(&entity.User{}).Where("id = ?", user.ID).Update("type", entity.AdminUser)
	if result.Error != nil {
		t.FailNow()
	}

	accessToken, refreshToken, _ = logUserIn(t, username, password)

	// make request to admin-only endpoint as admin
	_, qlRes = doQL(t, &accessToken, &refreshToken, query)

	// should not have errors
	assert.Nil(t, qlRes.Errors)
}
