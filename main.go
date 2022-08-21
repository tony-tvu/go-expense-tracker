package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/tony-tvu/goexpense/server"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found")
	}
	ctx := context.Background()
	s := server.Server{}
	s.Init(
		ctx,
		os.Getenv("PORT"),
		os.Getenv("ENV"),
		os.Getenv("SECRET"),
		os.Getenv("JWT_KEY"),
		os.Getenv("REFRESH_TOKEN_EXP"),
		os.Getenv("ACCESS_TOKEN_EXP"),
		os.Getenv("MONGODB_URI"),
		os.Getenv("DATABASE"),
		os.Getenv("PLAID_CLIENT_ID"),
		os.Getenv("PLAID_SECRET"),
		os.Getenv("PLAID_ENV"),
		os.Getenv("PLAID_COUNTRY_CODES"),
		os.Getenv("PLAID_PRODUCTS"),
	)

	// go func() {
	// 	log.Printf("Listening on port %s", s.App.Port)
	// 	err := s.App.Server.ListenAndServe()
	// 	if err != nil {
	// 		log.Println("While serving HTTP: ", err)
	// 	}
	// }()

	log.Printf("Listening on port %s", s.App.Port)
	err := s.App.Server.ListenAndServe()
	if err != nil {
		log.Println("While serving HTTP: ", err)
	}

}
