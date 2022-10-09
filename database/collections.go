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
	Enrollments  *mongo.Collection
	Sessions     *mongo.Collection
	Transactions *mongo.Collection
	Users        *mongo.Collection
}

func (db *MongoDb) SetCollections(client *mongo.Client, dbName string) {
	db.Accounts = client.Database(dbName).Collection("accounts")
	db.Configs = client.Database(dbName).Collection("configs")
	db.Enrollments = client.Database(dbName).Collection("enrollments")
	db.Sessions = client.Database(dbName).Collection("sessions")
	db.Transactions = client.Database(dbName).Collection("transactions")
	db.Users = client.Database(dbName).Collection("users")
}

func (db *MongoDb) CreateUniqueConstraints(ctx context.Context) {
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
	if _, err := db.Enrollments.Indexes().CreateOne(
		ctx, mongo.IndexModel{
			Keys:    bson.D{{Key: "access_token", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	); err != nil {
		log.Fatal(err)
	}
	if _, err := db.Accounts.Indexes().CreateOne(
		ctx, mongo.IndexModel{
			Keys:    bson.D{{Key: "account_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	); err != nil {
		log.Fatal(err)
	}
}

// Creates initial admin user. Account details can be specified in .env
func (db *MongoDb) CreateInitialAdminUser(ctx context.Context, username, email, password string) {
	// check if admin already exists
	count, err := db.Users.CountDocuments(ctx, bson.M{"username": username})
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
