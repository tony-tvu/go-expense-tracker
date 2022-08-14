package main

import (
	"context"
	"log"

	"github.com/tony-tvu/goexpense/app"
	"github.com/tony-tvu/goexpense/config"
	"github.com/tony-tvu/goexpense/database"
)

func main() {
	cfg, err := config.GetAppConfig()
	if err != nil {
		log.Fatal(err)
	}

	mongoclient, err := database.GetMongoClient()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := mongoclient.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	app := &app.App{
		Config:      cfg,
		MongoClient: mongoclient,
	}

	app.Run()
}
