package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/tony-tvu/goexpense/server"
)

func main() {
	ctx := context.Background()
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found")
	}

	s := server.Server{}
	mongoclient := s.Init(ctx,
		os.Getenv("ENV"),
		os.Getenv("PORT"),
		os.Getenv("AUTH_KEY"),
		os.Getenv("JWT_KEY"),
		os.Getenv("REFRESH_TOKEN_EXP"),
		os.Getenv("ACCESS_TOKEN_EXP"),
		os.Getenv("MONGODB_URI"),
		os.Getenv("DB_NAME"),
		os.Getenv("PLAID_CLIENT_ID"),
		os.Getenv("PLAID_SECRET"),
		os.Getenv("PLAID_ENV"),
		os.Getenv("PLAID_PRODUCTS"),
		os.Getenv("PLAID_COUNTRY_CODES"),
	)

	// must defer here to keep mongo connection alive
	defer func() {
		if err := mongoclient.Disconnect(ctx); err != nil {
			log.Println("mongo has been disconnected: ", err)
		}
	}()

	s.Start()
}
