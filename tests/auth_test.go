package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tony-tvu/goexpense/models"
	"go.mongodb.org/mongo-driver/bson"
)

func TestEmailLogin(t *testing.T) {
	t.Parallel()

	// given
	name := "TestEmailLoginName"
	email := "TestEmailLogin@email.com"
	password := "TestEmailLoginPassword"

	// create user
	m, b := map[string]string{
		"name":     name,
		"email":    email,
		"password": password},
		new(bytes.Buffer)
	json.NewEncoder(b).Encode(m)
	http.Post(fmt.Sprintf("%s/api/create_user", srv.URL), "application/json", b)

	// when: login with wrong password
	m, b = map[string]string{
		"email":    email,
		"password": "wrongPassword"},
		new(bytes.Buffer)
	json.NewEncoder(b).Encode(m)
	res, _ := http.Post(fmt.Sprintf("%s/api/login", srv.URL), "application/json", b)

	// then: 403 returned
	assert.Equal(t, http.StatusForbidden, res.StatusCode)

	// and: login with unfound email
	m, b = map[string]string{
		"email":    "unfound@email.com",
		"password": password},
		new(bytes.Buffer)
	json.NewEncoder(b).Encode(m)
	res, _ = http.Post(fmt.Sprintf("%s/api/login", srv.URL), "application/json", b)

	// then: 404 returned
	assert.Equal(t, http.StatusNotFound, res.StatusCode)

	// and: login with correct password
	m, b = map[string]string{
		"email":    email,
		"password": password},
		new(bytes.Buffer)
	json.NewEncoder(b).Encode(m)
	res, _ = http.Post(fmt.Sprintf("%s/api/login", srv.URL), "application/json", b)

	assert.Equal(t, http.StatusOK, res.StatusCode)

	// and: user login session should be created
	var u *models.User
	s.App.Users.FindOne(ctx, bson.D{{Key: "email", Value: email}}).Decode(&u)

	var ss *models.Session
	s.App.Sessions.FindOne(ctx, bson.D{{Key: "user_id", Value: u.ObjectID.Hex()}}).Decode(&ss)

	assert.NotNil(t, ss)
}
