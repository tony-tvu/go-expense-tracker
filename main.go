package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/playwright-community/playwright-go"
	"github.com/tony-tvu/goexpense/app"
)

func main() {
	err := playwright.Install()
	if err != nil {
		log.Fatal(err)
	}
	
	ctx := context.Background()
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found")
	}

	a := app.App{}
	mongoclient := a.Initialize(ctx,
		os.Getenv("ENV"),
		os.Getenv("PORT"),
		os.Getenv("AUTH_KEY"),
		os.Getenv("MONGODB_URI"),
		os.Getenv("DB_NAME"))

	// must defer here to keep mongo connection alive
	defer func() {
		if err := mongoclient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	a.Run()
}
