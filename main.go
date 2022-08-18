package main

import (
	// "context"
	// "log"
	// "os"

	// "github.com/joho/godotenv"
	// "github.com/tony-tvu/goexpense/app"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/tony-tvu/goexpense/jobs"
)

func main() {
	// ctx := context.Background()
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found")
	}

	// a := app.App{}
	// mongoclient := a.Init(ctx,
	// 	os.Getenv("ENV"),
	// 	os.Getenv("PORT"),
	// 	os.Getenv("AUTH_KEY"),
	// 	os.Getenv("MONGODB_URI"),
	// 	os.Getenv("DB_NAME"))

	// // must defer here to keep mongo connection alive
	// defer func() {
	// 	if err := mongoclient.Disconnect(ctx); err != nil {
	// 		panic(err)
	// 	}
	// }()

	// a.Run()

	pc := jobs.PlaidClient{}
	pc.Init(
		os.Getenv("PLAID_CLIENT_ID"),
		os.Getenv("PLAID_SECRET"),
		os.Getenv("PLAID_ENV"),
		os.Getenv("PLAID_PRODUCTS"),
		os.Getenv("PLAID_COUNTRY_CODES"))
}
