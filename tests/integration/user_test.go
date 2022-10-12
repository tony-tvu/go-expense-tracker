package tests

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/types"
	"github.com/tony-tvu/goexpense/user"
	"go.mongodb.org/mongo-driver/bson"
)

// Login and Logout handlers work correctly
func TestLoginAndLogout(t *testing.T) {
	t.Parallel()

	// create user
	testUser, cleanup := createTestUser(t)
	defer cleanup()

	// login with wrong password
	_, _, statusCode := logUserIn(t, testUser.Username, "wrong")

	// should return 403
	assert.Equal(t, http.StatusForbidden, statusCode)

	// login with unknown username
	_, _, statusCode = logUserIn(t, "userNameDoesntExist", testUser.Password)

	// should return 404
	assert.Equal(t, http.StatusNotFound, statusCode)

	// login with correct credentials
	body := map[string]string{
		"username": testUser.Username,
		"password": testUser.Password,
	}
	res := makeRequest(t, "POST", "/api/login", nil, nil, body)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// should have user session saved in db
	var s *auth.Session
	if err := testApp.Db.Sessions.FindOne(ctx, bson.M{"user_id": testUser.ID}).Decode(&s); err != nil {
		t.FailNow()
	}
	assert.Equal(t, testUser.ID, s.UserID)

	// logout
	cookies := getCookies(t, res.Cookies())
	accessToken := cookies["goexpense_access"]
	refreshToken := cookies["goexpense_refresh"]
	res = makeRequest(t, "POST", "/api/logout", &accessToken, &refreshToken)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// should no longer have user session saved after logging out
	count, err := testApp.Db.Sessions.CountDocuments(ctx, bson.M{"user_id": testUser.ID})
	if err != nil {
		t.FailNow()
	}
	assert.Equal(t, int64(0), count)
}

// UserInfo resolver returns correct information
func TestUserInfo(t *testing.T) {
	t.Parallel()

	// create user
	testUser, cleanup := createTestUser(t)
	defer cleanup()

	accessToken, refreshToken, _ := logUserIn(t, testUser.Username, testUser.Password)
	res := makeRequest(t, "GET", "/api/user_info", &accessToken, &refreshToken)

	// should return 200
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// should have correct user info returned
	var u *user.User
	json.NewDecoder(res.Body).Decode(&u)
	assert.Equal(t, testUser.Username, u.Username)
	assert.Equal(t, testUser.Email, u.Email)
	assert.Equal(t, "", u.Password)
	assert.Equal(t, types.RegularUser, u.UserType)
}

// IsLoggedIn route should return correct values
func TestIsAdminRoute(t *testing.T) {
	t.Parallel()

	type LoggedInResponse struct {
		IsLoggedIn bool `json:"logged_in"`
		IsAdmin    bool `json:"is_admin"`
	}

	// make request to logged_in route as a guest user
	res := makeRequest(t, "GET", "/api/logged_in", nil, nil)

	// should return 200, logged_in: false, is_admin:false
	var resBody *LoggedInResponse
	json.NewDecoder(res.Body).Decode(&resBody)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.False(t, resBody.IsLoggedIn)
	assert.False(t, resBody.IsAdmin)

	// create regular user and login
	user, cleanup := createTestUser(t)
	defer cleanup()
	accessToken, refreshToken, _ := logUserIn(t, user.Username, user.Password)

	// make request to logged_in route as regular user
	res = makeRequest(t, "GET", "/api/logged_in", &refreshToken, &accessToken)

	// should return 200, logged_in: true, is_admin:false
	json.NewDecoder(res.Body).Decode(&resBody)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.True(t, resBody.IsLoggedIn)
	assert.False(t, resBody.IsAdmin)

	// logout user
	res = makeRequest(t, "POST", "/api/logout", &accessToken, &refreshToken)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// make user an admin and login
	_, err := testApp.Db.Users.UpdateOne(
		ctx,
		bson.M{"username": user.Username},
		bson.M{
			"$set": bson.M{
				"user_type":  types.AdminUser,
				"updated_at": time.Now(),
			}},
	)
	if err != nil {
		t.FailNow()
	}
	accessToken, refreshToken, _ = logUserIn(t, user.Username, user.Password)

	// make request to logged_in endpoint as admin
	res = makeRequest(t, "GET", "/api/logged_in", &accessToken, &refreshToken)

	// should return 200, logged_in: true, is_admin: true
	json.NewDecoder(res.Body).Decode(&resBody)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.True(t, resBody.IsLoggedIn)
	assert.True(t, resBody.IsAdmin)
}
