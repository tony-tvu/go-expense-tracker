package smoketest

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tony-tvu/goexpense/models"
	"go.mongodb.org/mongo-driver/bson"
)

func TestUserHandlers(t *testing.T) {
	t.Run("Login handler verifies inputs and save new user session upon success", func(t *testing.T) {
		t.Parallel()

		// create user 
		name := "UserSesh1"
		email := "userSesh@email.com"
		password := "password"
		cleanup := createUser(t, testApp.Db, name, email, password)
		defer cleanup()

		// login with invalid email
		_, _, statusCode := logUserIn(t, "notAnEmail", password)

		// should return 400
		assert.Equal(t, http.StatusBadRequest, statusCode)

		// login with wrong password
		_, _, statusCode = logUserIn(t, email, "wrong")

		// should return 403
		assert.Equal(t, http.StatusForbidden, statusCode)

		// login with unknown email
		_, _, statusCode = logUserIn(t, "unknown@email.com", password)

		// should return 404
		assert.Equal(t, http.StatusNotFound, statusCode)

		// login with correct credentials
		body := map[string]string{
			"email":    email,
			"password": password,
		}
		res := MakeApiRequest(t, "POST", "/api/login", nil, nil, body)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		// should have user session saved in db
		count, _ := testApp.Db.Sessions.CountDocuments(ctx, bson.D{{Key: "email", Value: email}})
		assert.Equal(t, 1, int(count))

		// logout
		cookies := getCookies(t, res.Cookies())
		accessToken := cookies["goexpense_access"]
		refreshToken := cookies["goexpense_refresh"]
		res = MakeApiRequest(t, "POST", "/api/logout", &accessToken, &refreshToken)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		// should delete user's session from db
		count, _ = testApp.Db.Sessions.CountDocuments(ctx, bson.D{{Key: "email", Value: email}})
		assert.Equal(t, 0, int(count))
	})

	t.Run("GetUserInfo handler returns correct information", func(t *testing.T) {
		t.Parallel()

		// create user
		name := "GetUserInfo"
		email := "GetUserInfo@email.com"
		password := "^%#(GY%H=G$%asdf"
		cleanup := createUser(t, testApp.Db, name, email, password)
		defer cleanup()

		accessToken, refreshToken, _ := logUserIn(t, email, password)
		res := MakeApiRequest(t, "GET", "/api/user_info", &accessToken, &refreshToken)

		// should return 200
		assert.Equal(t, http.StatusOK, res.StatusCode)

		// should have correct user info returned
		var u *models.User
		json.NewDecoder(res.Body).Decode(&u)
		assert.Equal(t, name, u.Name)
		assert.Equal(t, email, u.Email)
		assert.Equal(t, "", u.Password)
		assert.Equal(t, models.RegularUser, u.Type)
	})

	t.Run("Logout handler should delete all sessions on success", func(t *testing.T) {
		t.Parallel()

		// create user and login
		name := "Logout"
		email := "Logout@email.com"
		password := "^%#(GY%H=G$%asdf"
		cleanup := createUser(t, testApp.Db, name, email, password)
		defer cleanup()
		accessToken, refreshToken, _ := logUserIn(t, email, password)

		// should have user session saved after logging in
		count, _ := testApp.Db.Sessions.CountDocuments(ctx, bson.D{{Key: "email", Value: email}})
		assert.Equal(t, 1, int(count))

		// logout
		res := MakeApiRequest(t, "POST", "/api/logout", &accessToken, &refreshToken)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		// should no longer have user session saved after logging out
		count, _ = testApp.Db.Sessions.CountDocuments(ctx, bson.D{{Key: "email", Value: email}})
		assert.Equal(t, 0, int(count))
	})
}
