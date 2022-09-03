package smoketest

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tony-tvu/goexpense/entity"
	"gorm.io/gorm"
)

func TestUserHandlers(t *testing.T) {
	t.Run("Login and Logout handlers work correctly", func(t *testing.T) {
		t.Parallel()

		// create user
		username := "UserSesh1"
		email := "userSesh@email.com"
		password := "password"
		cleanup := createUser(t, username, email, password)
		defer cleanup()

		// login with wrong password
		_, _, statusCode := logUserIn(t, username, "wrong")

		// should return 403
		assert.Equal(t, http.StatusForbidden, statusCode)

		// login with unknown username
		_, _, statusCode = logUserIn(t, "unknown", password)

		// should return 404
		assert.Equal(t, http.StatusNotFound, statusCode)

		// login with correct credentials
		body := map[string]string{
			"username": username,
			"password": password,
		}
		res := MakeApiRequest(t, "POST", "/api/login", nil, nil, body)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		// should have user session saved in db
		var s entity.Session
		result := testApp.Db.Where("username = ?", username).First(&s)
		assert.Nil(t, result.Error)
		assert.Equal(t, username, s.Username)

		// logout
		cookies := getCookies(t, res.Cookies())
		accessToken := cookies["goexpense_access"]
		refreshToken := cookies["goexpense_refresh"]
		res = MakeApiRequest(t, "POST", "/api/logout", &accessToken, &refreshToken)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		// should no longer have user session saved after logging out
		result = testApp.Db.Where("username = ?", username).First(&s)
		assert.Equal(t, gorm.ErrRecordNotFound, result.Error)
	})

	t.Run("GetUserInfo handler returns correct information", func(t *testing.T) {
		t.Parallel()

		// create user
		username := "GetUserInfo"
		email := "GetUserInfo@email.com"
		password := "^%#(GY%H=G$%asdf"
		cleanup := createUser(t, username, email, password)
		defer cleanup()

		accessToken, refreshToken, _ := logUserIn(t, username, password)
		res := MakeApiRequest(t, "GET", "/api/user_info", &accessToken, &refreshToken)

		// should return 200
		assert.Equal(t, http.StatusOK, res.StatusCode)

		// should have correct user info returned
		var u *entity.User
		json.NewDecoder(res.Body).Decode(&u)
		assert.Equal(t, username, u.Username)
		assert.Equal(t, email, u.Email)
		assert.Equal(t, "", u.Password)
		assert.Equal(t, entity.RegularUser, u.Type)
	})
}
