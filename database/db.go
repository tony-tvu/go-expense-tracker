package database

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Db struct {
	Users             *mongo.Collection
	Sessions          *mongo.Collection
}
