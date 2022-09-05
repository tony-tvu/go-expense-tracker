package main

import (
	"context"

	"github.com/tony-tvu/goexpense/app"
)

func main() {
	ctx := context.Background()
	app := &app.App{}
	app.Initialize(ctx)
	app.Serve()
}
