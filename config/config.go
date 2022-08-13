package config

import "go.mongodb.org/mongo-driver/mongo"

type Config struct {
	Client         *mongo.Client
	Database       string
	DBTimeout      int
	UserCollection string
	AuthKey        []byte
}
