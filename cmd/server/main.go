package main

import (
	"context"

	"github.com/avnovoselov/live_debugger/internal/server"
)

func main() {
	s := server.NewServer("1.0.0")
	s.Run(context.Background())
}
