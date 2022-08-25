package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tony-tvu/goexpense/auth"
	"github.com/tony-tvu/goexpense/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// LoggedIn routes should be accessible to logged in users
func TestLoggedIn(t *testing.T) {
	t.Parallel()

	// given
	name := "LoggedInName"
	email := "LoggedIn@email.com"
	password := "LoggedInPassword"

	// create user
	m, b := map[string]string{
		"name":     name,
		"email":    email,
		"password": password},
		new(bytes.Buffer)
	json.NewEncoder(b).Encode(m)
	http.Post(fmt.Sprintf("%s/api/create_user", srv.URL), "application/json", b)

	// when: user logs in
	client := &http.Client{}
	m, b = map[string]string{
		"email":    email,
		"password": password},
		new(bytes.Buffer)
	json.NewEncoder(b).Encode(m)

	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/login", srv.URL), b)
	res, _ := client.Do(req)

	// get cookies from response
	cookies := GetCookies(t, res.Cookies())

	// then: access_token cookie should exist in login response
	original_token := cookies["goexpense_access"]
	assert.NotEmpty(t, original_token)

	// and: make request to GetUserInfo endpoint with access_token
	req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/user_info", srv.URL), nil)
	req.AddCookie(&http.Cookie{
		Name:  "goexpense_access",
		Value: original_token})
	res, _ = client.Do(req)

	// then: status OK and access_token is unchanged
	assert.Equal(t, http.StatusOK, res.StatusCode)
	cookies = GetCookies(t, res.Cookies())
	assert.Equal(t, "", cookies["goexpense_access"])

	// and: wait for access_token to expire
	time.Sleep(1 * time.Second)
	// make request with expired token
	res, _ = client.Do(req)

	// then: status OK and access_token is refreshed
	assert.Equal(t, http.StatusOK, res.StatusCode)
	cookies = GetCookies(t, res.Cookies())
	assert.NotEqual(t, original_token, cookies["goexpense_access"])

	// and: returned user info is correct
	var u *user.User
	json.NewDecoder(res.Body).Decode(&u)
	assert.Equal(t, name, u.Name)
	assert.Equal(t, email, u.Email)
	assert.Equal(t, "", u.Password)
	assert.Equal(t, false, u.Verified)
	assert.Equal(t, user.ExternalUser, u.Role)

	// and: wait for refresh_token to expire
	time.Sleep(1 * time.Second)

	// make request with expired access_token
	res, _ = client.Do(req)

	// then: 401 unauthorized returns and no cookie set
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	cookies = GetCookies(t, res.Cookies())
	assert.Equal(t, "", cookies["goexpense_access"])
}

// LoggedIn routes should not be accessible when users aren't logged in
func TestLoggedIn2(t *testing.T) {
	t.Parallel()

	// given
	client := &http.Client{}

	// when: request made with no access token cookie
	res, _ := http.Get(fmt.Sprintf("%s/api/user_info", srv.URL))

	// then: 401 unauthorized returned
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

	// and: access_token cookie with missing claims set
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/user_info", srv.URL), nil)
	req.AddCookie(&http.Cookie{
		Name:  "goexpense_access",
		Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"})
	res, _ = client.Do(req)

	// then: 401 unauthorized returned
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

	// and: access_token cookie with correct claims that's expired
	req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/user_info", srv.URL), nil)
	req.AddCookie(&http.Cookie{
		Name:  "goexpense_access",
		Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNjMwNjJjMmQ1ZjMzOGQzNzBlY2VkZDE2IiwiUm9sZSI6IkV4dGVybmFsIiwiZXhwIjoxNjYxMzQ4OTEwLCJpYXQiOjE2NjEzNDg5MDl9.OhYTbTAsVvXHgnkzqsRXSvLLZYuAWJdTVUd-u-AsCD0"})
	res, _ = client.Do(req)

	// then: 401 unauthorized returned
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func TestAdmin(t *testing.T) {
	t.Parallel()

	// given
	name := "TestAdminName"
	email := "TestAdmin@email.com"
	password := "TestAdminPassword"
	client := &http.Client{}

	// create user
	m, b := map[string]string{
		"name":     name,
		"email":    email,
		"password": password},
		new(bytes.Buffer)
	json.NewEncoder(b).Encode(m)
	http.Post(fmt.Sprintf("%s/api/create_user", srv.URL), "application/json", b)

	// when: request made to GetSessions by non-admin user and logged out
	res, _ := http.Get(fmt.Sprintf("%s/api/sessions", srv.URL))

	// then: 401 unauthorized returned
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

	// and: log user in as non-admin
	m, b = map[string]string{
		"email":    email,
		"password": password},
		new(bytes.Buffer)
	json.NewEncoder(b).Encode(m)

	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/login", srv.URL), b)
	res, _ = client.Do(req)

	// get cookies from response
	cookies := GetCookies(t, res.Cookies())
	access_token := cookies["goexpense_access"]

	// make request to GetSessions
	req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/sessions", srv.URL), nil)
	req.AddCookie(&http.Cookie{
		Name:  "goexpense_access",
		Value: access_token})
	res, _ = client.Do(req)

	// then: 401 unauthorized returned
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

	// and: logout user
	req, _ = http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/logout", srv.URL), nil)
	req.AddCookie(&http.Cookie{
		Name:  "goexpense_access",
		Value: access_token})
	res, _ = client.Do(req)

	// then: logout is successful
	assert.Equal(t, http.StatusOK, res.StatusCode)
	var u *user.User
	s.App.Users.FindOne(
		ctx, bson.D{{Key: "email", Value: email}}).Decode(&u)
	assert.Equal(t, name, u.Name)

	var ss *auth.Session
	err := s.App.Sessions.FindOne(
		ctx, bson.D{{Key: "user_id", Value: u.ObjectID.Hex()}}).Decode(&ss)

	// sessions with current user should be deleted
	assert.Equal(t, mongo.ErrNoDocuments, err)

	// and: make user an admin
	s.App.Users.UpdateOne(
		ctx,
		bson.M{"email": email},
		bson.D{
			{Key: "$set", Value: bson.D{{Key: "role", Value: user.AdminUser}}},
		},
	)

	// log user in as admin
	m, b = map[string]string{
		"email":    email,
		"password": password},
		new(bytes.Buffer)
	json.NewEncoder(b).Encode(m)

	req, _ = http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/login", srv.URL), b)
	res, _ = client.Do(req)

	// get cookies from response
	cookies = GetCookies(t, res.Cookies())
	access_token = cookies["goexpense_access"]

	// make request to GetSessions as admin
	req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/sessions", srv.URL), nil)
	req.AddCookie(&http.Cookie{
		Name:  "goexpense_access",
		Value: access_token})
	res, _ = client.Do(req)

	// then: 200 success returned
	assert.Equal(t, http.StatusOK, res.StatusCode)
}
