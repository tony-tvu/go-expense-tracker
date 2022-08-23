package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/tony-tvu/goexpense/app"
	"github.com/tony-tvu/goexpense/server"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found")
	}
	ctx := context.Background()
	s := &server.Server{
		App: &app.App{
			Env:               os.Getenv("ENV"),
			Port:              os.Getenv("PORT"),
			Secret:            os.Getenv("SECRET"),
			JwtKey:            os.Getenv("JWT_KEY"),
			RefreshTokenExp:   86400,
			AccessTokenExp:    900,
			MongoURI:          os.Getenv("MONGODB_URI"),
			Db:                os.Getenv("DATABASE"),
			PlaidClientID:     os.Getenv("PLAID_CLIENT_ID"),
			PlaidSecret:       os.Getenv("PLAID_SECRET"),
			PlaidEnv:          os.Getenv("PLAID_ENV"),
			PlaidCountryCodes: "US,CA",
			PlaidProducts:     "auth,transactions",
		},
	}
	s.Initialize(ctx)
	s.Run(ctx)
}
