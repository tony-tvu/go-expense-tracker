package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tony-tvu/goexpense/models"
	"go.mongodb.org/mongo-driver/bson"
)

func TestAuthHandlers(t *testing.T) {

	t.Run("EmailLogin handler should work correctly and save new user session upon successful login", func(t *testing.T) {
		// create user
		createUser(t, testApp.Db, name, email, password)

		// login with invalid email
		_, statusCode := logUserIn(t, "notAnEmail", password)

		// should return 400
		assert.Equal(t, http.StatusBadRequest, statusCode)

		// login with wrong password
		_, statusCode = logUserIn(t, email, "wrong")

		// should return 403
		assert.Equal(t, http.StatusForbidden, statusCode)

		// login with unknown email
		_, statusCode = logUserIn(t, "unknown@email.com", password)

		// should return 404
		assert.Equal(t, http.StatusNotFound, statusCode)

		// login with correct password
		_, statusCode = logUserIn(t, email, password)

		assert.Equal(t, http.StatusOK, statusCode)

		// should have user session saved in db
		var ss *models.Session
		testApp.Db.Sessions.FindOne(ctx, bson.D{{Key: "email", Value: email}}).Decode(&ss)
		assert.NotNil(t, ss)
	})
}
