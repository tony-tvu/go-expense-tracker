package tests

import (
	"context"
	"testing"

	"github.com/tony-tvu/goexpense/models"
	"github.com/tony-tvu/goexpense/user"
	"go.mongodb.org/mongo-driver/bson"
)

// TODO:
func TestUserHandlers(t *testing.T) {

	t.Run("", func(t *testing.T) {
		// create user
		user.SaveUser(context.TODO(), testApp.Db, &models.User{
			Name:     name,
			Email:    email,
			Password: password,
		})

		// make user an admin
		testApp.Db.Users.UpdateOne(
			ctx,
			bson.M{"email": email},
			bson.D{
				{Key: "$set", Value: bson.D{{Key: "type", Value: models.AdminUser}}},
			},
		)

		// login and get cookies

		// when: invited with bad email
	})

}
