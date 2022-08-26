package tests

import (
	"context"
	"testing"

	"github.com/tony-tvu/goexpense/models"
	"github.com/tony-tvu/goexpense/user"
	"go.mongodb.org/mongo-driver/bson"
)

func TestInvite(t *testing.T) {
	t.Parallel()

	// given
	name := "TestInvite"
	email := "TestInvite@email.com"
	password := "TestInvitePassword"

	// create user
	user.SaveUser(context.TODO(), s.App, &models.User{
		Name:     name,
		Email:    email,
		Password: password,
	})

	// make user an admin
	s.App.Users.UpdateOne(
		ctx,
		bson.M{"email": email},
		bson.D{
			{Key: "$set", Value: bson.D{{Key: "type", Value: models.AdminUser}}},
		},
	)
	

}
