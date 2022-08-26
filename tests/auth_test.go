package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tony-tvu/goexpense/models"
	"github.com/tony-tvu/goexpense/user"
	"go.mongodb.org/mongo-driver/bson"
)

func TestEmailLogin(t *testing.T) {
	t.Parallel()

	// given
	name := "TestEmailLoginName"
	email := "TestEmailLogin@email.com"
	password := "TestEmailLoginPassword"

	// create user
	user.SaveUser(context.TODO(), s.App, &models.User{
		Name:     name,
		Email:    email,
		Password: password,
	})

	// when: login with wrong password
	m, b := map[string]string{
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
	var ss *models.Session
	s.App.Sessions.FindOne(ctx, bson.D{{Key: "email", Value: email}}).Decode(&ss)

	assert.NotNil(t, ss)
}
