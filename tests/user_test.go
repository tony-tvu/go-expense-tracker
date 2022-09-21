package tests

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tony-tvu/goexpense/models"
	"go.mongodb.org/mongo-driver/bson"
)

// Login and Logout handlers work correctly
func TestLoginAndLogout(t *testing.T) {
	t.Parallel()

	// create user
	user, cleanup := createTestUser(t)
	defer cleanup()

	// login with wrong password
	_, _, statusCode := logUserIn(t, user.Username, "wrong")

	// should return 403
	assert.Equal(t, http.StatusForbidden, statusCode)

	// login with unknown username
	_, _, statusCode = logUserIn(t, "userNameDoesntExist", user.Password)

	// should return 404
	assert.Equal(t, http.StatusNotFound, statusCode)

	// login with correct credentials
	body := map[string]string{
		"username": user.Username,
		"password": user.Password,
	}
	res := makeRequest(t, "POST", "/api/login", nil, nil, body)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// should have user session saved in db
	var s *models.Session
	if err := testApp.Db.Sessions.FindOne(ctx, bson.D{{Key: "user_id", Value: user.ID}}).Decode(&s); err != nil {
		t.FailNow()
	}
	assert.Equal(t, user.ID, s.UserID)

	// logout
	cookies := getCookies(t, res.Cookies())
	accessToken := cookies["goexpense_access"]
	refreshToken := cookies["goexpense_refresh"]
	res = makeRequest(t, "POST", "/api/logout", &accessToken, &refreshToken)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// should no longer have user session saved after logging out
	count, err := testApp.Db.Users.CountDocuments(ctx, bson.D{{Key: "user_id", Value: user.ID}})
	if err != nil {
		t.FailNow()
	}
	assert.Equal(t, 0, count)
}

// UserInfo resolver returns correct information
func TestUserInfo(t *testing.T) {
	t.Parallel()

	// create user
	user, cleanup := createTestUser(t)
	defer cleanup()

	accessToken, refreshToken, _ := logUserIn(t, user.Username, user.Password)
	res := makeRequest(t, "GET", "/api/user_info", &accessToken, &refreshToken)

	// should return 200
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// should have correct user info returned
	var u *models.User
	json.NewDecoder(res.Body).Decode(&u)
	assert.Equal(t, user.Username, u.Username)
	assert.Equal(t, user.Email, u.Email)
	assert.Equal(t, "", u.Password)
	assert.Equal(t, models.RegularUser, u.Type)
}

// IsAdmin route should return 200 for admin accounts only
func TestIsAdminRoute(t *testing.T) {
	t.Parallel()

	// make request to is_admin route as a guest user
	res := makeRequest(t, "GET", "/api/is_admin", nil, nil)

	// should return 401
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

	// create regular user and login
	user, cleanup := createTestUser(t)
	defer cleanup()
	accessToken, refreshToken, _ := logUserIn(t, user.Username, user.Password)

	// make request to is_admin route as regular user
	res = makeRequest(t, "GET", "/api/is_admin", &refreshToken, &accessToken)

	// should return 401
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

	// logout user
	res = makeRequest(t, "POST", "/api/logout", &accessToken, &refreshToken)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// make user an admin and login
	_, err := testApp.Db.Users.UpdateOne(
		ctx,
		bson.M{"username": user.Username},
		bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "type", Value: models.AdminUser},
				{Key: "updated_at", Value: time.Now()},
			}},
		},
	)
	if err != nil {
		t.FailNow()
	}
	accessToken, refreshToken, _ = logUserIn(t, user.Username, user.Password)

	// make request to is_admin endpoint as admin
	res = makeRequest(t, "GET", "/api/sessions", &accessToken, &refreshToken)

	// should return 200
	assert.Equal(t, http.StatusOK, res.StatusCode)
}
