package user

import (
	"context"
	"time"

	"github.com/tony-tvu/goexpense/database"
	"github.com/tony-tvu/goexpense/models"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// Saves new user to db once they've accepted email invitation
func SaveUser(ctx context.Context, db *database.Db, u *models.User) (*models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	doc := &bson.D{
		{Key: "email", Value: u.Email},
		{Key: "name", Value: u.Name},
		{Key: "password", Value: string(hash)},
		{Key: "type", Value: models.RegularUser},
		{Key: "created_at", Value: time.Now()},
	}
	_, err = db.Users.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}

	return u, nil
}
