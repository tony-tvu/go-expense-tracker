package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Env            string
	Database       string
	DBTimeout      int
	UserCollection string
	AuthKey        []byte
	Port           string
}

func GetAppConfig() (a *AppConfig, err error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found")
	}
	env := os.Getenv("ENV")
	mongoURI := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("DATABASE_NAME")
	authKeyStr := os.Getenv("KEY")
	port := os.Getenv("PORT")

	if env == "" || mongoURI == "" || dbName == "" || authKeyStr == "" || port == "" {
		return nil, errors.New("error - missing environment variables")
	}

	return &AppConfig{
		Env:            env,
		Database:       dbName,
		DBTimeout:      5,
		UserCollection: "users",
		AuthKey:        []byte(authKeyStr),
		Port:           port,
	}, nil
}

func GetTestAppConfig() (a *AppConfig, err error) {
	return &AppConfig{
		Env:            "TEST",
		Database:       "goexpense_test",
		DBTimeout:      5,
		UserCollection: "users",
		AuthKey:        []byte("TestKeyThatIs32CharactersLong011"),
		Port:           "8081",
	}, nil
}
