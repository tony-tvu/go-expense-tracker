package main

import (
	"context"
	"github.com/tony-tvu/goexpense/server"
)

func main() {
	ctx := context.Background()
	s := server.Server{}
	s.Init(ctx)
	s.Run(ctx)
}
