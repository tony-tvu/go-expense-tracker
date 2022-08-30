package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tony-tvu/goexpense/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestMiddlware(t *testing.T) {

	t.Run("AuthRequired middleware should issue access tokens correctly", func(t *testing.T) {
		name := "LoggedInName"
		email := "LoggedIn@email.com"
		password := "LoggedInPassword"

		// create user and login
		createUser(t, testApp.Db, name, email, password)
		accessToken, _ := logUserIn(t, email, password)

		// make request to endpoint where user must be logged in
		client := &http.Client{}
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/user_info", srv.URL), nil)
		req.AddCookie(&http.Cookie{
			Name:  "goexpense_access",
			Value: accessToken})
		res, _ := client.Do(req)

		// should return 200
		assert.Equal(t, http.StatusOK, res.StatusCode)

		// should not return a new access token
		cookies := getCookies(t, res.Cookies())
		assert.Equal(t, "", cookies["goexpense_access"])

		// wait for access token to expire
		time.Sleep(time.Duration(accessTokenExp) * time.Second)

		// make request with expired access token
		res, _ = client.Do(req)

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
		assert.Equal(t, false, u.Verified)
		assert.Equal(t, models.RegularUser, u.Type)

		// wait for refresh token to expire
		time.Sleep(2 * time.Second)

		// make request with expired refresh and access tokens
		res, _ = client.Do(req)

		// should return 401 unauthorized
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

		// response should not have new token cookie
		cookies = getCookies(t, res.Cookies())
		assert.Equal(t, "", cookies["goexpense_access"])

	})

	t.Run("AuthRequired middleware should not allow access to guest users", func(t *testing.T) {
		// make request to endpoint where user must be logged in without access token
		client := &http.Client{}
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/user_info", srv.URL), nil)
		res, _ := client.Do(req)

		// should return 401
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

		// make request to same endpoint with invalid access token
		req.AddCookie(&http.Cookie{
			Name:  "goexpense_access",
			Value: "invalid"})
		res, _ = client.Do(req)

		// should return 401
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

		// create user and login
		createUser(t, testApp.Db, name, email, password)
		accessToken, _ := logUserIn(t, email, password)

		// make request to same endpoint when logged in
		req.AddCookie(&http.Cookie{
			Name:  "goexpense_access",
			Value: accessToken})
		res, _ = client.Do(req)

		// should return 200
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	})

	t.Run("AdminRequired middleware should not allow access to guest or regular users", func(t *testing.T) {
		client := &http.Client{}

		// make request to admin-only route as a guest user
		res, _ := http.Get(fmt.Sprintf("%s/sessions", srv.URL))

		// should return 401
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

		// create user and login
		createUser(t, testApp.Db, name, email, password)
		access_token, _ := logUserIn(t, email, password)

		// make request to admin-only route as regulard user
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/sessions", srv.URL), nil)
		req.AddCookie(&http.Cookie{
			Name:  "goexpense_access",
			Value: access_token})
		res, _ = client.Do(req)

		// should return 401
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

		// logout user
		req, _ = http.NewRequest(http.MethodPost, fmt.Sprintf("%s/logout", srv.URL), nil)
		req.AddCookie(&http.Cookie{
			Name:  "goexpense_access",
			Value: access_token})
		res, _ = client.Do(req)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		// should no longer have user session saved after logging out
		var ss *models.Session
		err := testApp.Db.Sessions.FindOne(ctx, bson.D{{Key: "email", Value: email}}).Decode(&ss)
		assert.Equal(t, mongo.ErrNoDocuments, err)

		// make user an admin
		testApp.Db.Users.UpdateOne(
			ctx,
			bson.M{"email": email},
			bson.D{
				{Key: "$set", Value: bson.D{{Key: "type", Value: models.AdminUser}}},
			},
		)

		// login
		access_token, _ = logUserIn(t, email, password)

		// make request to admin-only endpoint as admin
		req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/sessions", srv.URL), nil)
		req.AddCookie(&http.Cookie{
			Name:  "goexpense_access",
			Value: access_token})
		res, _ = client.Do(req)

		// should return 200
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}
