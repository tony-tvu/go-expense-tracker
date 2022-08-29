package main

import (
	"context"

	"github.com/tony-tvu/goexpense/app"
	"github.com/tony-tvu/goexpense/database"
)

func main() {
	app := &app.App{
		Db: &database.Db{},
	}
	app.Run(context.Background())
}
