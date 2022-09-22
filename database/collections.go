package database

import (
	"context"
	"log"
	"time"

	"github.com/tony-tvu/goexpense/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type MongoDb struct {
	Accounts     *mongo.Collection
	Configs      *mongo.Collection
	Items        *mongo.Collection
	Sessions     *mongo.Collection
	Transactions *mongo.Collection
	Users        *mongo.Collection
}

func GetItems(ctx context.Context, db *MongoDb) ([]*models.Item, error) {
	cursor, err := db.Items.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var items []*models.Item
	if err = cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	return items, nil
}

func CreateUniqueConstraints(ctx context.Context, db *MongoDb) {
	if _, err := db.Users.Indexes().CreateOne(
		ctx, mongo.IndexModel{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	); err != nil {
		log.Fatal(err)
	}
	if _, err := db.Users.Indexes().CreateOne(
		ctx, mongo.IndexModel{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	); err != nil {
		log.Fatal(err)
	}
	if _, err := db.Transactions.Indexes().CreateOne(
		ctx, mongo.IndexModel{
			Keys:    bson.D{{Key: "transaction_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	); err != nil {
		log.Fatal(err)
	}
}

// Creates initial admin user. Account details can be specified in .env
func CreateInitialAdminUser(ctx context.Context, db *MongoDb, username, email, password string) {
	// check if admin already exists
	count, err := db.Users.CountDocuments(ctx, bson.D{{Key: "username", Value: username}})
	if err != nil {
		log.Fatal(err)
	}
	if count == 1 {
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	doc := &bson.D{
		{Key: "username", Value: username},
		{Key: "email", Value: email},
		{Key: "password", Value: string(hash)},
		{Key: "type", Value: models.AdminUser},
		{Key: "created_at", Value: time.Now()},
		{Key: "updated_at", Value: time.Now()},
	}
	if _, err = db.Users.InsertOne(ctx, doc); err != nil {
		log.Fatal(err)
	}
}
