package database

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDb struct {
	Configs      *mongo.Collection
	Items        *mongo.Collection
	Sessions     *mongo.Collection
	Transactions *mongo.Collection
	Users        *mongo.Collection
}
