package smoketest

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tony-tvu/goexpense/models"
	"go.mongodb.org/mongo-driver/bson"
)

func TestMiddlware(t *testing.T) {
	t.Run("AuthRequired middleware should issue access tokens correctly", func(t *testing.T) {
		t.Parallel()

		// create user and login
		name := "middleware"
		email := "middleware@email.com"
		password := "^%#(GY%H=G$%asdf"
		createUser(t, testApp.Db, name, email, password)
		accessToken, _ := logUserIn(t, email, password)

		// make request to endpoint where user must be logged in
		res := MakeApiRequest(t, "GET", "/api/user_info", nil, &accessToken)

		// should return 200
		assert.Equal(t, http.StatusOK, res.StatusCode)

		// should not return a new access token
		cookies := getCookies(t, res.Cookies())
		assert.Equal(t, "", cookies["goexpense_access"])

		// wait for access token to expire
		time.Sleep(time.Duration(accessTokenExp) * time.Second)

		// make request with expired access token
		res = MakeApiRequest(t, "GET", "/api/user_info", nil, &accessToken)

		// should return 200
		assert.Equal(t, http.StatusOK, res.StatusCode)

		// should have new access token in response
		cookies = getCookies(t, res.Cookies())
		assert.NotEqual(t, accessToken, cookies["goexpense_access"])

		// should have correct user info returned
		var u *models.User
		json.NewDecoder(res.Body).Decode(&u)
		assert.Equal(t, name, u.Name)
		assert.Equal(t, email, u.Email)
		assert.Equal(t, "", u.Password)
		assert.Equal(t, models.RegularUser, u.Type)

		// wait for refresh token to expire
		time.Sleep(2 * time.Second)

		// make request with expired refresh and access tokens
		res = MakeApiRequest(t, "GET", "/api/user_info", nil, &accessToken)

		// should return 401 unauthorized
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

		// response should not have new token cookie
		cookies = getCookies(t, res.Cookies())
		assert.Equal(t, "", cookies["goexpense_access"])

	})

	t.Run("AuthRequired middleware should not allow access to guest users", func(t *testing.T) {
		t.Parallel()

		// make request to endpoint where user must be logged in without access token
		res := MakeApiRequest(t, "GET", "/api/user_info", nil, nil)

		// should return 401
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

		// make request to same endpoint with invalid access token
		invalidToken := "invalidToken"
		res = MakeApiRequest(t, "GET", "/api/user_info", nil, &invalidToken)

		// should return 401
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

		// create user and login
		name := "middleware2"
		email := "middleware2@email.com"
		password := "^%#(GY%H=G$%asdf"
		createUser(t, testApp.Db, name, email, password)
		accessToken, _ := logUserIn(t, email, password)

		// make request to same endpoint when logged in
		res = MakeApiRequest(t, "GET", "/api/user_info", nil, &accessToken)

		// should return 200
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("AdminRequired middleware should not allow access to guest or regular users", func(t *testing.T) {
		t.Parallel()

		// make request to admin-only route as a guest user
		res := MakeApiRequest(t, "GET", "/api/sessions", nil, nil)

		// should return 401
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

		// create user and login
		name := "middleware3"
		email := "middleware3@email.com"
		password := "^%#(GY%H=G$%asdf"
		createUser(t, testApp.Db, name, email, password)
		accessToken, _ := logUserIn(t, email, password)

		// make request to admin-only route as regular user
		res = MakeApiRequest(t, "GET", "/api/sessions", nil, &accessToken)

		// should return 401
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

		// logout user
		res = MakeApiRequest(t, "POST", "/api/logout", nil, &accessToken)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		// should no longer have user session saved after logging out
		count, _ := testApp.Db.Sessions.CountDocuments(ctx, bson.D{{Key: "email", Value: email}})
		assert.Equal(t, 0, int(count))

		// make user an admin and login
		testApp.Db.Users.UpdateOne(
			ctx,
			bson.M{"email": email},
			bson.D{
				{Key: "$set", Value: bson.D{{Key: "type", Value: models.AdminUser}}},
			},
		)
		accessToken, _ = logUserIn(t, email, password)

		// make request to admin-only endpoint as admin
		res = MakeApiRequest(t, "GET", "/api/sessions", nil, &accessToken)

		// should return 200
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}
