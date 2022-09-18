package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tony-tvu/goexpense/entity"
	"gorm.io/gorm"
)

// Login and Logout handlers work correctly
func TestLoginAndLogout(t *testing.T) {
	t.Parallel()

	// create user
	username := "loginTestUser"
	email := "loginTestUser@email.com"
	password := "password"
	user, cleanup := createUser(t, username, email, password)
	defer cleanup()

	// login with wrong password
	_, _, statusCode := logUserIn(t, username, "wrong")

	// should return 403
	assert.Equal(t, http.StatusForbidden, statusCode)

	// login with unknown username
	_, _, statusCode = logUserIn(t, "userNameDoesntExist", password)

	// should return 404
	assert.Equal(t, http.StatusNotFound, statusCode)

	// login with correct credentials
	body := map[string]string{
		"username": username,
		"password": password,
	}
	res := makeRequest(t, "POST", "/api/login", nil, nil, body)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// should have user session saved in db
	var s entity.Session
	result := testApp.Db.Where("user_id = ?", user.ID).First(&s)
	assert.Nil(t, result.Error)
	assert.Equal(t, user.ID, s.UserID)

	// logout
	cookies := getCookies(t, res.Cookies())
	accessToken := cookies["goexpense_access"]
	refreshToken := cookies["goexpense_refresh"]
	res = makeRequest(t, "POST", "/api/logout", &accessToken, &refreshToken)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// should no longer have user session saved after logging out
	result = testApp.Db.Where("user_id = ?", user.ID).First(&s)
	assert.Equal(t, gorm.ErrRecordNotFound, result.Error)
}

// UserInfo resolver returns correct information
func TestUserInfo(t *testing.T) {
	t.Parallel()

	// create user
	username := "GetUserInfo"
	email := "GetUserInfo@email.com"
	password := "^%#(GY%H=G$%asdf"
	user, cleanup := createUser(t, username, email, password)
	defer cleanup()

	accessToken, refreshToken, _ := logUserIn(t, username, password)
	res := makeRequest(t, "GET", "/api/user_info", &accessToken, &refreshToken)

	// should return 200
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// should have correct user info returned
	var u *entity.User
	json.NewDecoder(res.Body).Decode(&u)
	assert.Equal(t, user.Username, u.Username)
	assert.Equal(t, user.Email, u.Email)
	assert.Equal(t, "", u.Password)
	assert.Equal(t, entity.RegularUser, u.Type)
}
