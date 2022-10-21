package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDb struct {
	Accounts     *mongo.Collection
	Enrollments  *mongo.Collection
	Rules        *mongo.Collection
	Sessions     *mongo.Collection
	Transactions *mongo.Collection
	Users        *mongo.Collection
}

func (db *MongoDb) SetCollections(client *mongo.Client, dbName string) {
	db.Accounts = client.Database(dbName).Collection("accounts")
	db.Enrollments = client.Database(dbName).Collection("enrollments")
	db.Rules = client.Database(dbName).Collection("rules")
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
