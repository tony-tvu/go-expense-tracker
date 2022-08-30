package main

import (
	"context"

	"github.com/tony-tvu/goexpense/app"
	"github.com/tony-tvu/goexpense/database"
)

func main() {
	ctx := context.Background()
	app := &app.App{
		Db: &database.Db{},
	}
	app.Initialize(ctx)
	app.Run(ctx)
}
