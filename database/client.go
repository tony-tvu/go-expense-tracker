package database

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMongoClient() (m *mongo.Client, err error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found")
	}
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return nil, errors.New("error - mongoURI is missing from environment")
	}

	mongoclient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	return mongoclient, nil
}

func GetMongoTestClient() (m *mongo.Client) {
	testDbURI := "mongodb://localhost:27017/local_db"
	mongoclient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testDbURI))
	if err != nil {
		panic(err)
	}

	return mongoclient
}
