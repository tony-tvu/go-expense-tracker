package tests

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tony-tvu/goexpense/entity"
	"gorm.io/gorm"
)

// Login and Logout resolvers work correctly
func TestLoginAndLogout(t *testing.T) {
	t.Parallel()

	// create user
	username := "loginTestUser"
	email := "loginTestUser@email.com"
	password := "password"
	user, cleanup := createUser(t, username, email, password)
	defer cleanup()

	// login with wrong password
	_, _, qlRes := logUserIn(t, username, "wrong")

	// should return error
	assert.Equal(t, "not authorized", qlRes.Errors[0].Message)

	// login with unknown username
	_, _, qlRes = logUserIn(t, "userNameDoesntExist", password)

	// should return unknown user error
	assert.Equal(t, "user not found", qlRes.Errors[0].Message)

	// login with correct credentials
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
	res, qlRes := doQL(t, nil, nil, query)

	// should not have errors
	assert.Nil(t, qlRes.Errors)

	// should have user session saved in db
	var s entity.Session
	result := testApp.Db.Where("user_id = ?", user.ID).First(&s)
	assert.Nil(t, result.Error)

	// logout
	cookies := getCookies(t, res.Cookies())
	accessToken := cookies["goexpense_access"]
	refreshToken := cookies["goexpense_refresh"]
	logoutQ :=
		`mutation {
		logout
	}`
	_, qlRes = doQL(t, &accessToken, &refreshToken, logoutQ)
	assert.Nil(t, qlRes.Errors)

	// should no longer have user session saved after logging out
	result = testApp.Db.Where("user_id = ?", user.ID).First(&s)
	assert.Equal(t, gorm.ErrRecordNotFound, result.Error)
}

// UserInfo resolver returns correct information
func TestUserInfo(t *testing.T) {
	t.Parallel()

	// create user
	username := "UserInfo"
	email := "UserInfo@email.com"
	password := "^%#(GY%H=G$%asdf"
	_, cleanup := createUser(t, username, email, password)
	defer cleanup()

	accessToken, refreshToken, _ := logUserIn(t, username, password)

	// make request
	query :=
		`query {
				userInfo {
					id
					username
					email
					type
					createdAt
					updatedAt
				}
			}`
	_, qlRes := doQL(t, &accessToken, &refreshToken, query)

	// should not have errors
	assert.Nil(t, qlRes.Errors)

	// should have correct user info returned
	var data struct {
		UserInfo entity.User
	}

	err := json.Unmarshal(qlRes.Data, &data)
	require.NoError(t, err)

	assert.NotNil(t, data.UserInfo.ID)
	assert.Equal(t, username, data.UserInfo.Username)
	assert.Equal(t, email, data.UserInfo.Email)
	assert.Equal(t, entity.RegularUser, data.UserInfo.Type)
	assert.NotNil(t, data.UserInfo.CreatedAt)
}