package smoketest

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tony-tvu/goexpense/user"
	"go.mongodb.org/mongo-driver/bson"
)

func TestUserHandlers(t *testing.T) {
	t.Run("handler verify inputs and save new user session upon successful login", func(t *testing.T) {
		t.Parallel()

		// create user
		name := "UserSesh"
		email := "userSesh@email.com"
		password := "^%#(GY%H=G$%asdf"
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
		body := user.CredentialsInput{
			Email: email,
			Password: password,
		}
		res := MakeApiRequest(t, "POST", "/login", &body, nil)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		// should have user session saved in db
		count, _ := testApp.Db.Sessions.CountDocuments(ctx, bson.D{{Key: "email", Value: email}})
		assert.Equal(t, 1, int(count))

		// logout
		cookies := getCookies(t, res.Cookies())
		accessToken := cookies["goexpense_access"]
		res = MakeApiRequest(t, "POST", "/logout", nil, &accessToken)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		// should delete user's session from db
		count, _ = testApp.Db.Sessions.CountDocuments(ctx, bson.D{{Key: "email", Value: email}})
		assert.Equal(t, 0, int(count))
	})

}
